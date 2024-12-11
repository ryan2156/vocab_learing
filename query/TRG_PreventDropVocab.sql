CREATE TRIGGER trg_PreventDropVocab
ON Vocabularies
AFTER DELETE
AS
BEGIN
    DECLARE @current_user_id INT
    SELECT @current_user_id = CONVERT(INT, CAST(CONTEXT_INFO() AS NVARCHAR))

    -- 如果嘗試刪除不是自己新增的單字，roll back操作
    IF EXISTS (
        SELECT 1
        FROM DELETED d
        WHERE d.added_by <> @current_user_id
    )
    BEGIN
        RAISERROR ('您無權刪除不是您新增的單字。', 16, 1)
        ROLLBACK TRANSACTION
    END
END