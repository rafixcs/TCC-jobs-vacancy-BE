-- goose Up
ALTER TABLE job_vacancies
ADD COLUMN search_vector tsvector;

UPDATE job_vacancies
SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, ''));

CREATE FUNCTION job_vacancies_search_vector_trigger() 
RETURNS trigger 
LANGUAGE plpgsql 
AS $func$ 
BEGIN 
  NEW.search_vector := 
    to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, '')); 
  RETURN NEW; 
END; 
$func$;

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON job_vacancies FOR EACH ROW EXECUTE FUNCTION job_vacancies_search_vector_trigger();

CREATE INDEX idx_job_vacancies_search_vector
ON job_vacancies USING GIN (search_vector);

-- +goose Down
-- DROP TRIGGER tsvectorupdate ON job_vacancies;
-- DROP FUNCTION job_vacancies_search_vector_trigger();
