CREATE TABLE `xkcd` (
  `xkcd_id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL,
  `url` varchar(255) NOT NULL,
  `file_content` blob DEFAULT NULL,    /* base64 format */
  `date_published` varchar(255) NOT NULL,
  `alt` text DEFAULT NULL,
  PRIMARY KEY (`xkcd_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin