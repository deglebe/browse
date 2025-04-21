# browse
world's worst browser engine

### features

##### pkg/html
- tag tokenization
    - start and end tags
    - skips `<!-- ... -->` comments
- attributes
    - double-quoted, single-quoted, and unquoted attributes
    - boolean attributes
- self-closing detection via `/` as well as minimal html5 void elements
- entity decoding
    - named entities
    - numeric decimal and hex entities
- void elements are treated as self-closing

##### dom

- `dom.Node`
    - `Type`, `Data`, `Attrs`, `Children`, `Parent`, `SelfClosing`
    - `PrettyPrint()` with sorted, quoted attributes and inline self-closing tags
- query helpers
    - `GetElementsByTagName(tag)`
    - `GetElementByID(id)`
    - `QuerySelectorAll("tag" | ".class" | "#id")`
    - `QuerySelector(...)`


