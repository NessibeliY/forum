sqlite3 forum.db < migrations/01_add_user_table.sql
sqlite3 forum.db < migrations/02_add_post_table.sql
sqlite3 forum.db < migrations/03_add_category_table.sql
sqlite3 forum.db < migrations/04_add_session_table.sql
sqlite3 forum.db < migrations/05_add_comment_table.sql
sqlite3 forum.db < migrations/06_add_post_reaction_table.sql
sqlite3 forum.db < migrations/07_add_comment_reaction_table.sql
# sqlite3 forum.db < migrations/08_insert_categories.sql
sqlite3 forum.db < migrations/09_add_post_category_table.sql
sqlite3 forum.db < migrations/10_add_image_table.sql
sqlite3 forum.db < migrations/11_add_notification_table.sql
sqlite3 forum.db < migrations/12_add_moderated_post_table.sql
sqlite3 forum.db < migrations/13_add_new_role_request_table.sql
# sqlite3 forum.db < migrations/14_alter_user_table.sql
sqlite3 forum.db < migrations/15_added_admin.sql
