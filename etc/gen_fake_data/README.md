# gen_fake_data

Generates fake data for use in testing.

## Steps to run

1. Create an empty database by starting the server with the --no-cron flag: e.g., `./monitor --no-cron`
2. Copy `monitor.db` from Step 1 into this folder
3. Run `Rscript gen-fake-data.R` to generate the fake data
4. Run `csvs-to-sqlite *.csv monitor.db` to add the data to `monitor.db`
5. Copy `monitor.db` back to the root of the project
6. Start the server with `./monitor --no-cron` again
