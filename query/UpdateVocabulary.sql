CREATE PROCEDURE UpdateVocabulary
    @VocabID INT,                 -- 單字的唯一 ID
    @Word NVARCHAR(70),           -- 單字
    @Defination NTEXT,            -- 定義
    @ExampleEng NTEXT = NULL,     -- 英文例句
    @ExampleZh NTEXT = NULL,      -- 中文例句
    @Part INT                     -- 詞性 ID
AS
BEGIN
    SET NOCOUNT ON;

    -- 確認 `VocabID` 是否存在於資料表中
    IF NOT EXISTS (SELECT 1 FROM Vocabularies WHERE vocab_id = @VocabID)
    BEGIN
        RAISERROR ('Vocabulary ID does not exist.', 16, 1);
        RETURN;
    END

    -- 更新語句
    UPDATE Vocabularies
    SET 
        word = @Word,
        defination = @Defination,
        example_eng = @ExampleEng,
        example_zh = @ExampleZh,
        part = @Part,
        added_date = GETDATE() -- 更新最後修改時間
    WHERE vocab_id = @VocabID;

    -- 返回成功訊息
    PRINT 'Vocabulary updated successfully.';
END;