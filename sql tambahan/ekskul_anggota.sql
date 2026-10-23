-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Sep 27, 2026 at 03:09 PM
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
-- Table structure for table `ekskul_anggota`
--

CREATE TABLE `ekskul_anggota` (
  `id_anggota` int NOT NULL,
  `nis` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `semester` enum('Ganjil','Genap') COLLATE utf8mb4_general_ci NOT NULL,
  `id_ekskul` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `ekskul_anggota`
--
ALTER TABLE `ekskul_anggota`
  ADD PRIMARY KEY (`id_anggota`),
  ADD UNIQUE KEY `uniq_ekskul_anggota` (`nis`,`id_tahunpelajaran`,`semester`,`id_ekskul`),
  ADD KEY `idx_ekskul_anggota_periode` (`id_tahunpelajaran`,`semester`,`id_ekskul`),
  ADD KEY `idx_ekskul_anggota_nis` (`nis`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `ekskul_anggota`
--
ALTER TABLE `ekskul_anggota`
  MODIFY `id_anggota` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
