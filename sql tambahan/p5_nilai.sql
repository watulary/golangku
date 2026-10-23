-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Sep 27, 2026 at 03:10 PM
-- Server version: 8.4.3
-- PHP Version: 8.3.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `smpn254`
--

-- --------------------------------------------------------

--
-- Table structure for table `p5_nilai`
--

CREATE TABLE `p5_nilai` (
  `id_p5` int NOT NULL,
  `nis` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `semester` enum('Ganjil','Genap') COLLATE utf8mb4_general_ci NOT NULL,
  `id_kelas` int NOT NULL,
  `dimensi` tinyint NOT NULL,
  `nilai` varchar(4) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `elemen` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sub_elemen` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `target` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `p5_nilai`
--
ALTER TABLE `p5_nilai`
  ADD PRIMARY KEY (`id_p5`),
  ADD UNIQUE KEY `uq_p5` (`nis`,`id_tahunpelajaran`,`semester`,`dimensi`),
  ADD KEY `idx_p5_nis` (`nis`),
  ADD KEY `fk_p5_tp` (`id_tahunpelajaran`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `p5_nilai`
--
ALTER TABLE `p5_nilai`
  MODIFY `id_p5` int NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `p5_nilai`
--
ALTER TABLE `p5_nilai`
  ADD CONSTRAINT `fk_p5_tp` FOREIGN KEY (`id_tahunpelajaran`) REFERENCES `tbl_tahunpelajaran` (`id_tahunpelajaran`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
