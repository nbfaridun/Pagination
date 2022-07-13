select * from sessions
offset ? limit ?;


CREATE PROCEDURE getMultipleSessions(index int)

    LANGUAGE plpgsql

AS $$

BEGIN;
-- Open a cursor for a query
DECLARE medley_cur CURSOR FOR SELECT * FROM sessions;
-- Retrieve ten rows
FETCH @index FROM medley_cur;
-- ...
-- Retrieve ten more from where we left off
FETCH 10 FROM medley_cur;
-- All done
COMMIT;
$$;


select * from sessions;

