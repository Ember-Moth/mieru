#!/bin/bash

# Copyright (C) 2024  mieru authors
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

function print_mieru_server_thread_dump() {
    echo "========== BEGIN OF MIERU SERVER THREAD DUMP =========="
    ./mita get thread-dump
    echo "==========  END OF MIERU SERVER THREAD DUMP  =========="
}

function print_mieru_server_metrics() {
    echo "========== BEGIN OF MIERU SERVER METRICS =========="
    ./mita get metrics
    ./mita get memory-statistics
    echo "==========  END OF MIERU SERVER METRICS  =========="
}
