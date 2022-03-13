# Generate fake data for ac-server-monitor
library(lubridate)
library(tibble)
library(readr)

sample_period <- dminutes(10)
dummy_time <- function() {
  now(tzone = "UTC")
}
status_id_counter <- 0

servers <- list(
  list(
    id = 1,
    name = "WintersEbb",
    begin = ymd_hms("2022-01-01T04:40:44Z"),
    end = ymd_hms("2022-01-02T04:40:44Z")
  ),
  list(
    id = 2,
    name = "Frostfell",
    begin = ymd_hms("2022-01-02T00:00:00Z"),
    end = ymd_hms("2022-01-03T00:00:00Z")
  )
)

gen_server_tbl <- function(id, name) {
  return (tibble(
    guid = id,
    name = name,
    description = paste0("Description for ", name),
    emu = "Emu",
    host = "localhost",
    port = "9000",
    type = "Type",
    status = "Status",
    website_url = "http://localhost",
    discord_url = "http://discord.com",
    is_listed = 1,
    created_at = as.integer(dummy_time()),
    updated_at = as.integer(dummy_time())
  ))
}

gen_status <- function(id, server_id, created_at) {
  tbl <- tibble(
    id = id,
    server_id = server_id,
    created_at = as.integer(created_at),
    status = rbinom(1, 1, 0.5)
  )
  
  return(tbl)
}

gen_statuses_tbl <- function(server_id, begin, end) {
  time_samples <- seq(begin, end, by = sample_period)
  
  do.call(rbind, lapply(time_samples, function(ts) {
    tbl <- gen_status(status_id_counter, server_id, ts)
    status_id_counter <<- status_id_counter + 1
    tbl
  }))
}

servers_tbl <- do.call(rbind, lapply(servers, function(x) {
  gen_server_tbl(x$id, x$name)
}))

statuses_tbl <- do.call(rbind, lapply(servers, function(x) {
  gen_statuses_tbl(x$id, x$begin, x$end)
}))

write_csv(servers_tbl, "servers.csv")
write_csv(statuses_tbl, "statuses.csv")