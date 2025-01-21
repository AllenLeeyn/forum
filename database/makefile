# Variables
DB_FILE = forum.db
SCHEMA_FILE = schema.txt
ACCOUNT_TYPE_FILE = accountType.txt
CATEGORY_FILE = category.txt
VIEW_FILE = view.txt

# Default target
all: $(DB_FILE)

# Create SQLite database
$(DB_FILE): $(SCHEMA_FILE) $(ACCOUNT_TYPE_FILE) $(CATEGORY_FILE) $(VIEW_FILE)
	@echo "Creating SQLite database..."
	# Create database and apply schema
	sqlite3 $@ < $(SCHEMA_FILE)
	# Insert account types
	@echo "Inserting account types..."
	sqlite3 $@ < $(ACCOUNT_TYPE_FILE)
	# Insert categories
	@echo "Inserting categories..."
	sqlite3 $@ < $(CATEGORY_FILE)
	# Create view
	@echo "Creating view..."
	sqlite3 $@ < $(VIEW_FILE)

# Clean target to remove the database file
clean:
	@echo "Cleaning up..."
	rm -f $(DB_FILE)
