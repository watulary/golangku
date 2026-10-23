-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Sep 27, 2026 at 03:32 PM
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
-- Table structure for table `tbl_tahunpelajaran`
--

CREATE TABLE `tbl_tahunpelajaran` (
  `id_tahunpelajaran` int NOT NULL,
  `nama_tahunpelajaran` varchar(9) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `semester` enum('Ganjil','Genap') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Ganjil',
  `aktif` int NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `tbl_tahunpelajaran`
--

INSERT INTO `tbl_tahunpelajaran` (`id_tahunpelajaran`, `nama_tahunpelajaran`, `semester`, `aktif`) VALUES
(1, '2012/2013', 'Genap', 0),
(2, '2013/2014', 'Genap', 0),
(3, '2014/2015', 'Genap', 0),
(4, '2015/2016', 'Genap', 0),
(5, '2016/2017', 'Genap', 0),
(6, '2017/2018', 'Genap', 0),
(7, '2018/2019', 'Genap', 0),
(8, '2019/2020', 'Genap', 0),
(9, '2020/2021', 'Genap', 0),
(10, '2021/2022', 'Genap', 0),
(11, '2022/2023', 'Genap', 0),
(12, '2023/2024', 'Genap', 0),
(13, '2024/2025', 'Genap', 0),
(14, '2025/2026', 'Genap', 0),
(15, '2026/2027', 'Ganjil', 1);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `tbl_tahunpelajaran`
--
ALTER TABLE `tbl_tahunpelajaran`
  ADD PRIMARY KEY (`id_tahunpelajaran`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `tbl_tahunpelajaran`
--
ALTER TABLE `tbl_tahunpelajaran`
  MODIFY `id_tahunpelajaran` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
