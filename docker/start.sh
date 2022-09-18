#!/bin/sh

/app/migrator -no-gen # generate/migrate tables
/app/server
