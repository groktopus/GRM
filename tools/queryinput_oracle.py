"""Regenerate shape differential cases from the checked-in normative schema.
Run with tools/oracle-requirements.txt. Leap applicability is a separate Go domain
oracle: this Python format checker rejects RFC3339 leap seconds.
"""
import json
from pathlib import Path
from jsonschema import Draft202012Validator, FormatChecker

root = Path(__file__).resolve().parents[1]
text = (root / 'docs/architecture/mcp-schemas.md').read_text()
registry = json.loads(text.split('```json\n', 1)[1].split('```', 1)[0])
schema = {'$schema': registry['$schema'], '$defs': registry['$defs'], '$ref': '#/$defs/QueryInput'}
Draft202012Validator.check_schema(schema)
checker = FormatChecker()
assert not checker.conforms('not-a-date', 'date-time')
assert not checker.conforms('2016-12-31T23:59:60Z', 'date-time')
validator = Draft202012Validator(schema, format_checker=checker)
cases = []
def add(value):
    cases.append({'input': json.dumps(value, ensure_ascii=False, separators=(',', ':')), 'valid': validator.is_valid(value)})
for branch in ['question', 'cursor']:
    base = {'contract': '1', branch: 'q'}
    add(base)
    for field in base:
        v = dict(base); del v[field]; add(v)
        for wrong in [None, False, 1, [], {}, '', '2']:
            add(dict(base, **{field: wrong}))
    for field in ['bindings', 'page_size', 'valid_at', 'known_at', 'unknown', 'Question']:
        for value in [None, False, 1, [], {}, 'x']:
            add(dict(base, **{field: value}))
for field in ['mention', 'entity']:
    base = {'mention': 'm', 'entity': 'ent_A'}
    values = [None, False, 1, [], {}, '']
    for value in values:
        add({'contract': '1', 'question': 'q', 'bindings': [dict(base, **{field: value})]})
    del base[field]
    add({'contract': '1', 'question': 'q', 'bindings': [base]})
for n in [0,1,20,21]:
    add({'contract':'1','question':'q','bindings':[{'mention':'m','entity':'ent_A'}]*n})
for value in [0,1,20,21,1.5]:
    add({'contract':'1','question':'q','page_size':value})
for value in ['2000-02-29T00:00:00Z','1900-02-29T00:00:00Z','2016-12-31T24:00:00Z','2016-12-31T23:59:59.0Z','2016-12-31T23:59:59+00:00']:
    for field in ['valid_at','known_at']:
        add({'contract':'1','question':'q',field:value})
for value in [{}, [], None, 1, True, 'x']:
    add(value)
path = root / 'internal/queryinput/testdata/schema-cases.json'
path.write_text(json.dumps(cases, ensure_ascii=False, indent=2)+'\n')
print(f'{len(cases)} cases; {sum(c["valid"] for c in cases)} accepted; format checking enabled; leap seconds handled separately')
