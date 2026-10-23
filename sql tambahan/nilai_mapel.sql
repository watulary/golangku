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
-- Table structure for table `nilai_mapel`
--

CREATE TABLE `nilai_mapel` (
  `id_nilai` int NOT NULL,
  `nis` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `semester` enum('Ganjil','Genap') COLLATE utf8mb4_general_ci NOT NULL,
  `id_kelas` int NOT NULL,
  `id_mapel` int NOT NULL,
  `nilai_akhir` decimal(5,2) DEFAULT NULL,
  `capaian_pembelajaran` text COLLATE utf8mb4_general_ci,
  `kktp` text COLLATE utf8mb4_general_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `nilai_mapel`
--
ALTER TABLE `nilai_mapel`
  ADD PRIMARY KEY (`id_nilai`),
  ADD UNIQUE KEY `uq_nilai` (`nis`,`id_tahunpelajaran`,`semester`,`id_mapel`),
  ADD KEY `idx_tp_sem` (`id_tahunpelajaran`,`semester`),
  ADD KEY `idx_mapel` (`id_mapel`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `nilai_mapel`
--
ALTER TABLE `nilai_mapel`
  MODIFY `id_nilai` int NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `nilai_mapel`
--
ALTER TABLE `nilai_mapel`
  ADD CONSTRAINT `fk_nilai_mapel` FOREIGN KEY (`id_mapel`) REFERENCES `tbl_mapel` (`id_mapel`),
  ADD CONSTRAINT `fk_nilai_tp` FOREIGN KEY (`id_tahunpelajaran`) REFERENCES `tbl_tahunpelajaran` (`id_tahunpelajaran`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
