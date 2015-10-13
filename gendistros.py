#!/usr/bin/env python
# Generate Go maps of distros supported by Packagecloud API

import os
import urllib
import json

def gen_map(name, distros):
    print 'var %s = map[string]int{' % name
    for d in distros:
        for v in d['versions']:
            k = '/'.join([d['index_name'], v['index_name']])
            v = v['id']
            print "\t\"%s\": %d," % (k, v)
    print '}'

token = os.environ['PACKAGECLOUD_TOKEN']
url = 'https://%s:@packagecloud.io/api/v1/distributions.json' % token
resp = urllib.urlopen(url)
data = json.loads(resp.read())

print '// Auto-generated with %s' % __file__
print
print 'package pkgcloud'
print
gen_map('debDistroIDs', data['deb'])
print
gen_map('rpmDistroIDs', data['rpm'])
