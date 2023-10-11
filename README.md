# uust-to-calendar
UUST-To-Calendar Parser.
1. Скачать репозиторий
2. Разархивировать
3. Открыть папку в терминале ![image](https://github.com/noodypv/uust-to-calendar/assets/126050017/061580d3-9dd8-4559-b80e-7f12f666746c)
4. Выполнить команду (пример ссылки - https://isu.ugatu.su/api/new_schedule_api/?schedule_semestr_id=231&WhatShow=1&student_group_id=2356&weeks=2)
   ```bash
   wt main.exe -u "ссылка на расписание вашей группы(любой недели)"
   ```
   ![image](https://github.com/noodypv/uust-to-calendar/assets/126050017/6ab68b9a-884c-4912-b320-3a01e8e31c08)

По истечении работы появится файл calendar.ics в той же директории где и файл main.exe. Если расписание изменилось, просто запустите парсер еще раз.

