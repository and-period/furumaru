CREATE SCHEMA IF NOT EXISTS `media` DEFAULT CHARACTER SET utf8mb4;

-- ライブ配信情報
CREATE TABLE IF NOT EXISTS `media`.`broadcasts` (
  `id`                           VARCHAR(22)  NOT NULL,          -- ライブ配信ID
  `schedule_id`                  VARCHAR(22)  NULL DEFAULT NULL, -- 開催スケジュールID
  `type`                         INT          NOT NULL,          -- ライブ配信種別
  `status`                       INT          NOT NULL,          -- ライブ配信状況
  `input_url`                    TEXT         NOT NULL,          -- ライブ配信URL(入力)
  `output_url`                   TEXT         NOT NULL,          -- ライブ配信URL(出力)
  `archive_url`                  TEXT         NOT NULL,          -- アーカイブ配信URL
  `cloud_front_distribution_arn` TEXT         NULL DEFAULT NULL, -- CloudFrontディストリビューションARN
  `media_live_channel_arn`       TEXT         NULL DEFAULT NULL, -- MediaLiveチャンネルARN
  `media_live_channel_id`        VARCHAR(256) NULL DEFAULT NULL, -- MediaLiveチャンネルID
  `media_live_rtmp_input_arn`    TEXT         NULL DEFAULT NULL, -- MediaLiveインプットARN(RTMP)
  `media_live_rtmp_input_name`   VARCHAR(256) NULL DEFAULT NULL, -- MediaLiveインプット名(RTMP)
  `media_live_mp4_input_arn`     TEXT         NULL DEFAULT NULL, -- MediaLiveインプットARN(MP4)
  `media_live_mp4_input_name`    VARCHAR(256) NULL DEFAULT NULL, -- MediaLiveインプット名(MP4)
  `media_store_container_arn`    TEXT         NULL DEFAULT NULL, -- MediaStoreコンテナARN
  `created_at`                   DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at`                   DATETIME(3)  NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_broadcast_schedule_id` ON `media`.`broadcasts` (`schedule_id` ASC) VISIBLE;
