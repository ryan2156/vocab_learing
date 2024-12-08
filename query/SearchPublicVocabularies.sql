CREATE PROCEDURE SearchPublicVocabularies
    @keyword NVARCHAR(100)
AS
BEGIN
    SELECT 
        vocab_id,
        word,
        defination,
        part_name,
        added_by_name,
        added_date
    FROM VW_PublicVocabs
    WHERE word LIKE '%' + @keyword + '%'
       OR defination LIKE '%' + @keyword + '%'
       OR part_name LIKE '%' + @keyword + '%'
    ORDER BY added_date DESC; -- 按加入日期排序
END;