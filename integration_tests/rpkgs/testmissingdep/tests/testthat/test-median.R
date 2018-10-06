context("test-median.R")

test_that("median works as expected", {
  vals <- c(1, 2, 3, 4, 5, 6, 7)
  expect_equal(my_median(vals), median(vals))
})

test_that("median shows up as failure", {
  vals <- c(1, 2, 3, 4, 5, 6, 7)
  # this is a failure we want
  expect_equal(my_median(vals), 0)
})
