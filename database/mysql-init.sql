DROP DATABASE IF EXISTS `komp_registry`;
CREATE DATABASE `komp_registry`;

USE `komp_registry`;

CREATE TABLE `Komps` (
  `id` binary(36) NOT NULL,
  `serial_number` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `state` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `software_version` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `product_code` enum('ev2','ev2b') NOT NULL,
  `mac_address` varchar(24) DEFAULT NULL,
  `comment` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_SERIAL_NUMBER` (`serial_number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `Attributes` (
  `komp_id` binary(36) NOT NULL,
  `name` varchar(64) NOT NULL,
  `value` varchar(64) NOT NULL,
  PRIMARY KEY (`komp_id`,`name`),
  CONSTRAINT `Attributes_ibfk_1` FOREIGN KEY (`komp_id`) REFERENCES `Komps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `Komps` VALUES
  ("d3fd6449-1b2d-4f6b-ad56-2f164a4ec0ec", "1puc6vaq", "allocated", "v3.3.3", "ev2b", "f8:dc:7a:08:a1:ce", "Allocated to BoVel/SYE on February 23rd 2023"),
  ("f226e9ff-99d6-4726-b558-da81231744de", "3pts8ff8", "allocated", "v3.3.0", "ev2b", "f8:dc:7a:08:a2:07", "At Oslo Met for testing/feature discovery/possibility scoping"),
  ("9794670e-0c6d-42be-b659-8bda0e87c7b2", "h8xerba5", "allocated", "v3.3.1", "ev2", "f8:dc:7a:1d:eb:b9", "Delivered to Østensjø last Friday (the 9th). So far so good!"),
  ("d355dc46-20c9-408d-94a7-5c2dd73ad6b3", "uqpfnn53", "available", "v3.3.2", "ev2", "f8:dc:7a:0d:eb:b8", ""),
  ("c0c7793e-4999-4340-9ed6-64f2ec052045", "q5ew966u", "available", "v3.3.2", "ev2b", "fa:dc:7a:08:a2:cc", "Expected to be delivered to Sognsvann, not yet allocated, awaiting final go ahead"),
  ("08de34e5-058e-4490-9b8a-cb7b465ae411", "7f9mekf4", "available", "v3.3.2", "ev2b", "f8:dc:7a:0d:eb:b4", ""),
  ("20793d81-6b57-4924-a930-5cd81e93531e", "52cmzym6", "available", "v3.3.2", "ev2b", "f8:dc:7a:08:a1:cf", ""),
  ("72715292-354d-4e12-b3e6-204239feef27", "t3nfpaxa", "available", "v3.3.2", "ev2", "f8:dc:7a:0d:fb:d3", "Test unit at Majorstua, needs to be allocated depending on results. TBC..."),
  ("2da109a5-0b0e-4df8-adf1-07991de235cc", "q3hgdjnp", "available", "v3.3.1", "ev2", "f8:dc:7a:0d:eb:d2", ""),
  ("85b52077-a034-42b8-939e-f352bd68a976", "fzhnfqd9", "available", "v3.2.1", "ev2", "f8:dc:7a:08:a1:05", ""),
  ("70eb11f7-d366-472e-a227-90c417534303", "h5zeraa5", "available", "v3.2.1", "ev2", "f8:dc:7a:08:a3:05", ""),
  ("7a6d489e-81b4-450d-b704-6821aba2cf0b", "zurwvy6y", "available", "v3.2.1", "ev2", "f8:dc:7a:0d:eb:b9", "");


INSERT INTO `Attributes` VALUES
  ("d3fd6449-1b2d-4f6b-ad56-2f164a4ec0ec", "simcard_iccid", "73450401200514377995"),
  ("d3fd6449-1b2d-4f6b-ad56-2f164a4ec0ec", "simcard_state", "active"),
  ("f226e9ff-99d6-4726-b558-da81231744de", "simcard_iccid", "83250331200514377127"),
  ("f226e9ff-99d6-4726-b558-da81231744de", "simcard_state", "blocked"),
  ("c0c7793e-4999-4340-9ed6-64f2ec052045", "simcard_iccid", "83250331200514831923"),
  ("c0c7793e-4999-4340-9ed6-64f2ec052045", "simcard_state", "inactive"),
  ("08de34e5-058e-4490-9b8a-cb7b465ae411", "simcard_iccid", "83250331200514921311"),
  ("08de34e5-058e-4490-9b8a-cb7b465ae411", "simcard_state", "active"),
  ("20793d81-6b57-4924-a930-5cd81e93531e", "simcard_iccid", "83250331200514913821"),
  ("20793d81-6b57-4924-a930-5cd81e93531e", "simcard_state", "inactive");