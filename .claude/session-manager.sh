#!/bin/bash

# Claude Code用セッション管理ユーティリティ
SESSION_FILE=".claude/session.json"
SESSION_BACKUP_DIR=".claude/debug/sessions"

# セッション更新関数
update_session() {
    local field=$1
    local value=$2
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # jqを使用してセッションを更新
    if command -v jq &> /dev/null; then
        jq --arg field "$field" --arg value "$value" --arg ts "$timestamp" \
            '.current_session[$field] = $value | .current_session.last_updated = $ts' \
            "$SESSION_FILE" > "$SESSION_FILE.tmp" && mv "$SESSION_FILE.tmp" "$SESSION_FILE"
    else
        echo "セッション管理にはjqが必要です"
        exit 1
    fi
}

# タスク追加関数
add_task() {
    local task=$1
    local status=${2:-"pending"}
    
    if command -v jq &> /dev/null; then
        jq --arg task "$task" --arg status "$status" \
            '.current_session.tasks[$status] += [$task]' \
            "$SESSION_FILE" > "$SESSION_FILE.tmp" && mv "$SESSION_FILE.tmp" "$SESSION_FILE"
    fi
}

# メモ追加関数
add_note() {
    local note=$1
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    if command -v jq &> /dev/null; then
        jq --arg note "$note" --arg ts "$timestamp" \
            '.current_session.notes += [{"timestamp": $ts, "content": $note}]' \
            "$SESSION_FILE" > "$SESSION_FILE.tmp" && mv "$SESSION_FILE.tmp" "$SESSION_FILE"
    fi
}

# セッションを履歴に保存する関数
save_session() {
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local backup_file="$SESSION_BACKUP_DIR/session_$timestamp.json"
    
    # バックアップディレクトリが存在しない場合は作成
    mkdir -p "$SESSION_BACKUP_DIR"
    
    # 現在のセッションをバックアップにコピー
    cp "$SESSION_FILE" "$backup_file"
    
    # 現在のセッションを履歴にアーカイブ
    if command -v jq &> /dev/null; then
        current_session=$(jq '.current_session' "$SESSION_FILE")
        jq --argjson session "$current_session" \
            '.history = [$session] + .history | .history = .history[0:100]' \
            "$SESSION_FILE" > "$SESSION_FILE.tmp" && mv "$SESSION_FILE.tmp" "$SESSION_FILE"
        
        # 現在のセッションをリセット
        new_session_id="session_$(date +"%Y%m%d_%H%M%S")"
        jq --arg id "$new_session_id" --arg ts "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
            '.current_session = {"id": $id, "started_at": $ts, "status": "active", "context": {}, "tasks": {"current": "", "completed": [], "pending": []}, "notes": []}' \
            "$SESSION_FILE" > "$SESSION_FILE.tmp" && mv "$SESSION_FILE.tmp" "$SESSION_FILE"
    fi
    
    echo "セッションを保存しました: $backup_file"
}

# 最近のセッション一覧表示関数
list_sessions() {
    if [ -d "$SESSION_BACKUP_DIR" ]; then
        echo "最近のセッション:"
        ls -1t "$SESSION_BACKUP_DIR" | head -10
    fi
    
    if command -v jq &> /dev/null && [ -f "$SESSION_FILE" ]; then
        echo -e "\n現在のセッション:"
        jq -r '.current_session | "ID: \(.id)\n開始: \(.started_at)\nステータス: \(.status)"' "$SESSION_FILE"
    fi
}

# メインコマンドハンドラー
case "$1" in
    update)
        update_session "$2" "$3"
        ;;
    task)
        add_task "$2" "$3"
        ;;
    note)
        add_note "$2"
        ;;
    save)
        save_session
        ;;
    list)
        list_sessions
        ;;
    *)
        echo "使用方法: $0 {update|task|note|save|list} [引数...]"
        echo "  update <フィールド> <値>     - セッションフィールドを更新"
        echo "  task <説明> [ステータス]     - タスクを追加 (ステータス: pending/completed)"
        echo "  note <内容>                  - メモを追加"
        echo "  save                         - 現在のセッションを保存して新規開始"
        echo "  list                         - 最近のセッション一覧"
        exit 1
        ;;
esac