import csv
import numpy as np
from scipy import stats

def read_csv(file_name):
    """Read a single line of integers from a CSV file."""
    with open(file_name, 'r') as file:
        reader = csv.reader(file)
        for row in reader:
            # Convert each value in the row to an integer
            return [int(x) for x in row]

def analyze_data(data):
    """Calculate statistical metrics for the list of integers."""
    mean = np.mean(data)
    variance = np.var(data, ddof=1)  # Use ddof=1 for sample variance
    std_dev = np.std(data, ddof=1)   # Use ddof=1 for sample standard deviation
    conf_interval = stats.t.interval(0.95, len(data)-1, loc=mean, scale=stats.sem(data))
    
    metrics = {
        'runs': len(data),
        'units': "miliseconds",
        'mean': mean,
        'variance': variance,
        'std_dev': std_dev,
        'conf_interval': conf_interval,
        'min': np.min(data),
        'max': np.max(data),
        'median': np.median(data),
        '25_percentile': np.percentile(data, 25),
        '75_percentile': np.percentile(data, 75),
    }
    
    return metrics

def write_metrics_to_file(metrics, file_name):
    """Write the calculated metrics to a text file."""
    with open(file_name, 'w') as file:
        for key, value in metrics.items():
            file.write(f'{key}: {value}\n')

def main():
    input_file = 'output.csv'   # Change to your input CSV file
    output_file = 'report.txt' # Change to your desired output file
    
    # Read data from CSV
    data = read_csv(input_file)
    
    # Analyze the data
    metrics = analyze_data(data)
    
    # Write metrics to the output file
    write_metrics_to_file(metrics, output_file)
    print(f'Metrics written to {output_file}')

if __name__ == "__main__":
    main()
