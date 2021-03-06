import os
import sys
import datetime

import crawler

if('MONTH' in os.environ):
    month = os.environ['MONTH']
else:
    sys.stderr.write("Invalid arguments, missing parameter: 'MONTH'.\n")
    os._exit(1)
if('YEAR' in os.environ):
    year = os.environ['YEAR']
else:
    sys.stderr.write("Invalid arguments, missing parameter: 'YEAR'.\n")
    os._exit(1)
if('DRIVER_PATH' in os.environ):
    driver_path = os.environ['DRIVER_PATH']
else:
    sys.stderr.write("Invalid arguments, missing parameter: 'DRIVER_PATH'.\n")
    os._exit(1)
if('OUTPUT_PATH' in os.environ):
    output_path = os.environ['OUTPUT_PATH']
else:
    output_path = "/output"

now = datetime.datetime.now()
current_year = now.year
current_month = now.month
if((int(month) < 1) | (int(month) > 12)):
    sys.stderr.write("Invalid month {}: InvalidParameters.\n".format(month))
    os._exit(1)
if((int(year) == current_year) & (int(month) > current_month)):
    sys.stderr.write("Invalid month {}: InvalidParameters.\n".format(month))
    os._exit(1)
if(int(year) > current_year):
    sys.stderr.write("Invalid year {}: InvalidParameters.\n".format(year))
    os._exit(1)

files = crawler.crawl(output_path, driver_path, month, year)
