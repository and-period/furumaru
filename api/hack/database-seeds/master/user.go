package master

import "github.com/and-period/furumaru/api/internal/user/entity"

var DummyUsers = []*entity.NewUserParams{
	{UserType: entity.UserTypeGuest, Lastname: "佐藤", Firstname: "優子", LastnameKana: "サトウ", FirstnameKana: "ユウコ", Email: "sato.yuko@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "鈴木", Firstname: "大翔", LastnameKana: "スズキ", FirstnameKana: "ダイオウ", Email: "suzuki.daiou@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "田中", Firstname: "美咲", LastnameKana: "タナカ", FirstnameKana: "ミサキ", Email: "tanaka.misaki@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "渡辺", Firstname: "翔太", LastnameKana: "ワタナベ", FirstnameKana: "ショウタ", Email: "watanabe.shota@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "伊藤", Firstname: "愛", LastnameKana: "イトウ", FirstnameKana: "アイ", Email: "ito.ai@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "中村", Firstname: "陽一", LastnameKana: "ナカムラ", FirstnameKana: "ヨウイチ", Email: "nakamura.yoichi@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "小林", Firstname: "凛", LastnameKana: "コバヤシ", FirstnameKana: "リン", Email: "kobayashi.rin@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "加藤", Firstname: "大和", LastnameKana: "カトウ", FirstnameKana: "ヤマト", Email: "kato.yamato@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "木村", Firstname: "結衣", LastnameKana: "キムラ", FirstnameKana: "ユイ", Email: "kimura.yui@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "山本", Firstname: "蓮", LastnameKana: "ヤマモト", FirstnameKana: "レン", Email: "yamamoto.ren@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "高橋", Firstname: "さくら", LastnameKana: "タカハシ", FirstnameKana: "サクラ", Email: "takahashi.sakura@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "斎藤", Firstname: "健一", LastnameKana: "サイトウ", FirstnameKana: "ケンイチ", Email: "saito.kenichi@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "松本", Firstname: "美咲", LastnameKana: "マツモト", FirstnameKana: "ミサキ", Email: "matsumoto.misaki@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "井上", Firstname: "大輝", LastnameKana: "イノウエ", FirstnameKana: "ダイキ", Email: "inoue.daiki@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "山田", Firstname: "千尋", LastnameKana: "ヤマダ", FirstnameKana: "チヒロ", Email: "yamada.chihiro@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "清水", Firstname: "悠斗", LastnameKana: "シミズ", FirstnameKana: "ユウト", Email: "shimizu.yuto@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "石川", Firstname: "美羽", LastnameKana: "イシカワ", FirstnameKana: "ミウ", Email: "ishikawa.miu@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "近藤", Firstname: "颯太", LastnameKana: "コンドウ", FirstnameKana: "ソウタ", Email: "kondo.sota@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "坂本", Firstname: "七海", LastnameKana: "サカモト", FirstnameKana: "ナナミ", Email: "sakamoto.nanami@example.com"},
	{UserType: entity.UserTypeGuest, Lastname: "吉田", Firstname: "翔太", LastnameKana: "ヨシダ", FirstnameKana: "ショウタ", Email: "yoshida.shota@example.com"},
}
