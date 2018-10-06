context("test-median.R")

test_that("median works as expected", {
  vals <- c(1, 2, 3, 4, 5, 6, 7)
  expect_equal(my_super_median(vals), stats::median(vals))
})
