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
-- Table structure for table `nilai_ekskul`
--

CREATE TABLE `nilai_ekskul` (
  `id_nilai_ekskul` int NOT NULL,
  `nis` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `semester` enum('Ganjil','Genap') COLLATE utf8mb4_general_ci NOT NULL,
  `id_kelas` int NOT NULL,
  `id_ekskul` int NOT NULL,
  `nilai` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `keterangan` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `nilai_ekskul`
--
ALTER TABLE `nilai_ekskul`
  ADD PRIMARY KEY (`id_nilai_ekskul`),
  ADD UNIQUE KEY `uq_nilai_ekskul` (`nis`,`id_tahunpelajaran`,`semester`,`id_ekskul`),
  ADD KEY `idx_tp_sem` (`id_tahunpelajaran`,`semester`),
  ADD KEY `idx_ekskul` (`id_ekskul`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `nilai_ekskul`
--
ALTER TABLE `nilai_ekskul`
  MODIFY `id_nilai_ekskul` int NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `nilai_ekskul`
--
ALTER TABLE `nilai_ekskul`
  ADD CONSTRAINT `fk_nilaiekskul_ekskul` FOREIGN KEY (`id_ekskul`) REFERENCES `tbl_ekskul` (`id_ekskul`),
  ADD CONSTRAINT `fk_nilaiekskul_tp` FOREIGN KEY (`id_tahunpelajaran`) REFERENCES `tbl_tahunpelajaran` (`id_tahunpelajaran`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
