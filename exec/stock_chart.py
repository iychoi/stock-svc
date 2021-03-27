#! /usr/bin/python3
import os
import os.path
import sys
import pandas as pd
import yfinance as yf
import matplotlib.pyplot as plt
#https://pypi.org/project/yfinance/
#https://www.codementor.io/@hachimy15/quantitative-finance-and-data-visualization-in-python-for-beginners-16apvwc49b
#https://matplotlib.org/2.0.2/examples/pylab_examples/simple_plot.html


def getData(ticker, period, interval):
    data = yf.download(ticker, period=period, interval=interval, progress=False)
    return data

def saveChart(data, period, filepath):
    data["Adj Close"].plot()
    plt.xlabel("Time %s" % period)
    plt.ylabel("Price")
    plt.savefig(filepath, bbox_inches='tight')

def main(argv):
    if len(argv) < 4:
        print("command : ./stock_chart.py ticker period interval filepath")
    else:
        ticker = argv[0]
        period = argv[1]
        interval = argv[2]
        filepath = argv[3]

        data = getData(ticker, period, interval)
        saveChart(data, period, filepath)

if __name__ == "__main__":
    main(sys.argv[1:])
