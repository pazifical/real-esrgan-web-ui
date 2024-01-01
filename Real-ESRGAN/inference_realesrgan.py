import argparse
import shutil

print("Simulating the processor just for testing without ERSGAN")

parser = argparse.ArgumentParser()
parser.add_argument("-i", "--input")
parser.add_argument("-o", "--output")
parser.add_argument("-s", "--scale")
parser.add_argument("-n")
args = parser.parse_args()

print(args.input, args.output, args.scale)

shutil.copy(args.input, args.output)
