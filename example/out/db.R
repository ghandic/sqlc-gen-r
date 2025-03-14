# Generate R functions for SQL queries
# Each function will:
# 1. Take parameters (if any), with type checks.
# 2. Execute the SQL query.
# 3. Return the results as a data frame.

library(DBI)
library(glue)
library(rlang)
source("R/import-standalone-obj-type.R")
source("R/import-standalone-types-check.R")

# Comments from SQL:

#' getReplyIds
#'
#' Executes the SQL query "getReplyIds" and returns the result as a data frame.
#'
#' @param conn A DBI connection object.
#' @param parent_id integer()
#' @return A data frame containing the results of the query.
#' @export
getReplyIds <- function(conn, parent_id) {
  # Type checks for parameters using r-lib's standalone checks
  check_number_whole(parent_id)

  # Build the SQL query using glue
  sql <- glue::glue_sql("SELECT
    id
FROM
    post
WHERE
    parent_id = ?",
    parent_id = parent_id,
    .con = conn
  )

  # Execute the query and fetch the results
  result <- DBI::dbGetQuery(conn, sql)

  return(result)
}
