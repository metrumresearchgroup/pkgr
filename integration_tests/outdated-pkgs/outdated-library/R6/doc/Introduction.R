## ----echo = FALSE--------------------------------------------------------
knitr::opts_chunk$set(collapse = TRUE, comment = "#>")

## ------------------------------------------------------------------------
library(R6)

Person <- R6Class("Person",
  public = list(
    name = NA,
    hair = NA,
    initialize = function(name, hair) {
      if (!missing(name)) self$name <- name
      if (!missing(hair)) self$hair <- hair
      self$greet()
    },
    set_hair = function(val) {
      self$hair <- val
    },
    greet = function() {
      cat(paste0("Hello, my name is ", self$name, ".\n"))
    }
  )
)

## ------------------------------------------------------------------------
ann <- Person$new("Ann", "black")
ann

## ------------------------------------------------------------------------
ann$hair
ann$greet()
ann$set_hair("red")
ann$hair

## ------------------------------------------------------------------------
Queue <- R6Class("Queue",
  public = list(
    initialize = function(...) {
      for (item in list(...)) {
        self$add(item)
      }
    },
    add = function(x) {
      private$queue <- c(private$queue, list(x))
      invisible(self)
    },
    remove = function() {
      if (private$length() == 0) return(NULL)
      # Can use private$queue for explicit access
      head <- private$queue[[1]]
      private$queue <- private$queue[-1]
      head
    }
  ),
  private = list(
    queue = list(),
    length = function() base::length(private$queue)
  )
)

q <- Queue$new(5, 6, "foo")

## ------------------------------------------------------------------------
# Add and remove items
q$add("something")
q$add("another thing")
q$add(17)
q$remove()
q$remove()

## ----eval = FALSE--------------------------------------------------------
#  q$queue
#  #> NULL
#  q$length()
#  #> Error: attempt to apply non-function

## ------------------------------------------------------------------------
q$add(10)$add(11)$add(12)

## ------------------------------------------------------------------------
q$remove()
q$remove()
q$remove()
q$remove()

## ------------------------------------------------------------------------
Numbers <- R6Class("Numbers",
  public = list(
    x = 100
  ),
  active = list(
    x2 = function(value) {
      if (missing(value)) return(self$x * 2)
      else self$x <- value/2
    },
    rand = function() rnorm(1)
  )
)

n <- Numbers$new()
n$x

## ------------------------------------------------------------------------
n$x2

## ------------------------------------------------------------------------
n$x2 <- 1000
n$x

## ----eval=FALSE----------------------------------------------------------
#  n$rand
#  #> [1] 0.2648
#  n$rand
#  #> [1] 2.171
#  n$rand <- 3
#  #> Error: unused argument (quote(3))

## ------------------------------------------------------------------------
# Note that this isn't very efficient - it's just for illustrating inheritance.
HistoryQueue <- R6Class("HistoryQueue",
  inherit = Queue,
  public = list(
    show = function() {
      cat("Next item is at index", private$head_idx + 1, "\n")
      for (i in seq_along(private$queue)) {
        cat(i, ": ", private$queue[[i]], "\n", sep = "")
      }
    },
    remove = function() {
      if (private$length() - private$head_idx == 0) return(NULL)
      private$head_idx <<- private$head_idx + 1
      private$queue[[private$head_idx]]
    }
  ),
  private = list(
    head_idx = 0
  )
)

hq <- HistoryQueue$new(5, 6, "foo")
hq$show()
hq$remove()
hq$show()
hq$remove()

## ------------------------------------------------------------------------
CountingQueue <- R6Class("CountingQueue",
  inherit = Queue,
  public = list(
    add = function(x) {
      private$total <<- private$total + 1
      super$add(x)
    },
    get_total = function() private$total
  ),
  private = list(
    total = 0
  )
)

cq <- CountingQueue$new("x", "y")
cq$get_total()
cq$add("z")
cq$remove()
cq$remove()
cq$get_total()

## ------------------------------------------------------------------------
SimpleClass <- R6Class("SimpleClass",
  public = list(x = NULL)
)

SharedField <- R6Class("SharedField",
  public = list(
    e = SimpleClass$new()
  )
)

s1 <- SharedField$new()
s1$e$x <- 1

s2 <- SharedField$new()
s2$e$x <- 2

# Changing s2$e$x has changed the value of s1$e$x
s1$e$x

## ------------------------------------------------------------------------
NonSharedField <- R6Class("NonSharedField",
  public = list(
    e = NULL,
    initialize = function() e <<- SimpleClass$new()
  )
)

n1 <- NonSharedField$new()
n1$e$x <- 1

n2 <- NonSharedField$new()
n2$e$x <- 2

# n2$e$x does not affect n1$e$x
n1$e$x

## ------------------------------------------------------------------------
RC <- setRefClass("RC",
  fields = list(x = 'ANY'),
  methods = list(
    getx = function() x,
    setx = function(value) x <<- value
  )
)

rc <- RC$new()
rc$setx(10)
rc$getx()

## ------------------------------------------------------------------------
NP <- R6Class("NP",
  portable = FALSE,
  public = list(
    x = NA,
    getx = function() x,
    setx = function(value) x <<- value
  )
)

np <- NP$new()
np$setx(10)
np$getx()

## ------------------------------------------------------------------------
P <- R6Class("P",
  portable = TRUE,  # This is default
  public = list(
    x = NA,
    getx = function() self$x,
    setx = function(value) self$x <- value
  )
)

p <- P$new()
p$setx(10)
p$getx()

## ------------------------------------------------------------------------
Simple <- R6Class("Simple",
  public = list(
    x = 1,
    getx = function() x
  )
)

Simple$set("public", "getx2", function() self$x*2)

# To replace an existing member, use overwrite=TRUE
Simple$set("public", "x", 10, overwrite = TRUE)

s <- Simple$new()
s$x
s$getx2()

## ------------------------------------------------------------------------
PrettyCountingQueue <- R6Class("PrettyCountingQueue",
  inherit = CountingQueue,
  public = list(
    print = function(...) {
      cat("<PrettyCountingQueue> of ", self$get_total(), " elements\n", sep = "")
      invisible(self)
    }
  )
)

## ------------------------------------------------------------------------
pq <- PrettyCountingQueue$new(1, 2, "foobar")
pq

