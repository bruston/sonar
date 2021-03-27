sonar
=====

Subdomain enumeration via the JSON API provided by https://sonar.omnisint.io/ which uses the Rapid7 dataset.

## Usage

`sonar -d example.com | tee -a example.com.hosts`

Does not save to a file by default. Instead, redirect the output to a file using > or tee, or pipe straight into another program.
