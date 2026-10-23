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
-- Table structure for table `nilai_ujian`
--

CREATE TABLE `nilai_ujian` (
  `id_ujian` int NOT NULL,
  `nis` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `jenis` enum('Sekolah','Daerah') COLLATE utf8mb4_general_ci NOT NULL,
  `id_kelas` int NOT NULL,
  `id_mapel` int NOT NULL,
  `nilai_akhir` decimal(5,2) DEFAULT NULL,
  `capaian_pembelajaran` text COLLATE utf8mb4_general_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `nilai_ujian`
--
ALTER TABLE `nilai_ujian`
  ADD PRIMARY KEY (`id_ujian`),
  ADD UNIQUE KEY `uq_ujian` (`nis`,`id_tahunpelajaran`,`jenis`,`id_mapel`),
  ADD KEY `idx_ujian_nis` (`nis`),
  ADD KEY `fk_ujian_tp` (`id_tahunpelajaran`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `nilai_ujian`
--
ALTER TABLE `nilai_ujian`
  MODIFY `id_ujian` int NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `nilai_ujian`
--
ALTER TABLE `nilai_ujian`
  ADD CONSTRAINT `fk_ujian_tp` FOREIGN KEY (`id_tahunpelajaran`) REFERENCES `tbl_tahunpelajaran` (`id_tahunpelajaran`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
