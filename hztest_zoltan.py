#!/usr/bin/env python3
import argparse
import re
import io
def accepted_levels(arg_level):
    levels = {
        'TRACE': ['TRACE', 'DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'],
        'DEBUG': ['DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'],
        'INFO': ['INFO', 'WARN', 'ERROR', 'FATAL'],
        'WARN': ['WARN', 'ERROR', 'FATAL'],
        'ERROR': ['ERROR', 'FATAL'],
        'FATAL': ['FATAL'],
    }
    return levels.get(arg_level)
def main():
    parser = argparse.ArgumentParser(description='Extract log from a specific Hazelcast test from a test output file.')
    parser.add_argument('--file', metavar='file', help='the file to read')
    parser.add_argument('--test', metavar='test', help='the id of the test to extract')
    parser.add_argument('--log-level', metavar='log_level', help='the minimum log level to show', default="TRACE", nargs='?')
    parser.add_argument('--log-class', metavar='log_class', help='the log lines to filter by comma separated list of logger classes', nargs='?')
    parser.add_argument('--log-member', metavar='log_member', help='the log lines to filter by member', nargs='?')
    args = parser.parse_args()
    re_log_prefix_pattern = re.compile('^(\d{2}:\d{2}:\d{2},\d{3})\s+(\w+).*\|\s+\-\s+\[(\w+)\]\s+([^\s]+)\s+\-\s+([^\s]+).*')

    test_log_file = open(args.file, 'r')
    matching_log_entry = False
    list_accepted_classes = [] if args.log_class is None else args.log_class.split(',')
    list_accepted_levels = accepted_levels(args.log_level)
    list_accepted_members = [] if args.log_member is None else args.log_member.split(',')
    while True:
        line = test_log_file.readline()
        if not line:
            break
        prefix_match = re_log_prefix_pattern.match(line)
        matching_test = args.test in line
        if prefix_match:
            # ts = prefix_match.group(1)
            log_level = prefix_match.group(2)
            log_class = prefix_match.group(3)
            # log_thread = prefix_match.group(4)
            log_member = prefix_match.group(5)
            accepted_level = any(log_level in s for s in list_accepted_levels)
            matching_class = len(list_accepted_classes) == 0 or any(log_class in s for s in list_accepted_classes)
            matching_member = len(list_accepted_members) == 0 or any(log_member in s for s in list_accepted_members)
            if not matching_test or not accepted_level or not matching_class or not matching_member:
                matching_log_entry = False
                continue
        if args.test in line:
            matching_log_entry = True
        if matching_log_entry:
            print(line.rstrip())
    test_log_file.close()
if __name__ == '__main__':
    main()