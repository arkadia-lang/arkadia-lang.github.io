-- Adapts docs/guide.md for the website: the sidebar shows the table of
-- contents, so the one written at the top of the guide (the first list and
-- the rule after it) is dropped, and the title moves from the body (where
-- the table of contents would nest every section under it) to the
-- `heading` variable of the template.
function Pandoc(doc)
  local blocks = pandoc.Blocks({})
  local list, rule = false, false
  for _, b in ipairs(doc.blocks) do
    if b.t == "BulletList" and not list then
      list = true
    elseif b.t == "HorizontalRule" and list and not rule then
      rule = true
    elseif b.t == "Header" and b.level == 1 and not doc.meta.heading then
      doc.meta.heading = b.content
    else
      blocks:insert(b)
    end
  end
  doc.blocks = blocks
  return doc
end
