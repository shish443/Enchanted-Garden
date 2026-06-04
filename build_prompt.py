import os

# ЕСЛИ СКРИПТ ВНУТРИ ПАПКИ ПРОЕКТА, ПОСТАВЬ ТОЧКУ: project_dir = '.'
project_dir = '/home/admin/Рабочий стол/библиотека моей жизни/03_Source_Code/goProjects/Enchanted-Garden' 
output_file = 'codebase_for_review.md'

ignore_dirs = {'.git', '.idea', '.vscode', 'vendor', 'bin', 'pkg'}
allowed_extensions = {'.go', '.md', '.json', '.mod', '.sum', '.yaml', '.yml'}

# Проверка, существует ли папка вообще
if not os.path.exists(project_dir):
    print(f"❌ ОШИБКА: Папка '{project_dir}' не найдена! Проверь путь.")
    exit()

files_added = 0

with open(output_file, 'w', encoding='utf-8') as outfile:
    for root, dirs, files in os.walk(project_dir):
        dirs[:] = [d for d in dirs if d not in ignore_dirs]

        for file in files:
            ext = os.path.splitext(file)[1].lower()
            if ext not in allowed_extensions and file not in ['Makefile', 'Dockerfile']:
                continue

            # Чтобы скрипт не читал сам себя или прошлые результаты
            if file == output_file or file == 'build_prompt.py':
                continue

            filepath = os.path.join(root, file)
            
            if ext == '.go': lang = 'go'
            elif ext == '.json': lang = 'json'
            elif ext in ['.yaml', '.yml']: lang = 'yaml'
            else: lang = ''

            try:
                with open(filepath, 'r', encoding='utf-8') as infile:
                    content = infile.read()

                outfile.write(f"### Файл: `{filepath}`\n\n```{lang}\n{content}\n```\n\n")
                files_added += 1
                print(f"Добавлен: {filepath}") # Пишем в консоль каждый найденный файл

            except Exception as e:
                print(f"⚠️ Пропуск файла {filepath}: {e}")

if files_added == 0:
    print(f"\n❌ Файл пустой. В папке '{project_dir}' не найдено ни одного файла с нужным расширением (.go, .md и т.д.).")
else:
    print(f"\n✅ Готово! Собрано файлов: {files_added}. Результат в {output_file}")