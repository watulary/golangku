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
-- Table structure for table `keputusan_akhir`
--

CREATE TABLE `keputusan_akhir` (
  `id_keputusan` int NOT NULL,
  `nis` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `id_tahunpelajaran` int NOT NULL,
  `status` enum('Lulus','Tidak Lulus') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `no_ijazah` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `no_skhus` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `tanggal` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `keputusan_akhir`
--
ALTER TABLE `keputusan_akhir`
  ADD PRIMARY KEY (`id_keputusan`),
  ADD UNIQUE KEY `uq_keputusan_nis` (`nis`),
  ADD KEY `fk_keputusan_tp` (`id_tahunpelajaran`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `keputusan_akhir`
--
ALTER TABLE `keputusan_akhir`
  MODIFY `id_keputusan` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `keputusan_akhir`
--
ALTER TABLE `keputusan_akhir`
  ADD CONSTRAINT `fk_keputusan_tp` FOREIGN KEY (`id_tahunpelajaran`) REFERENCES `tbl_tahunpelajaran` (`id_tahunpelajaran`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
