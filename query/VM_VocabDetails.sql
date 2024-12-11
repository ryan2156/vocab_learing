CREATE VIEW VocabDetails AS
SELECT 
    v.vocab_id, 
    v.word, 
    v.defination AS definition, 
    v.example_eng, 
    v.example_zh, 
    p.name AS part_name, 
    u.name AS added_by_name, 
    v.added_date
FROM 
    Vocabularies v
JOIN 
    Parts p ON v.part = p.part_id
JOIN 
    Users u ON v.added_by = u.user_id;