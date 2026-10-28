Контакты

Gmail: vladpurchel40@gmail.com
Telegram: Vladislav345687
Телефон: +7 909 650 10 19

📊 Расширенная модель угрозы

В соответствии с требованиями, стандартная модель была дополнена важными атрибутами контекста безопасности:
• md5 (string, 32 знака) — Уникальный хэш-идентификатор файла.
• name (string) — Сигнатурное имя угрозы (например, Trojan.Win32.Generic).
• danger_type (string) — Классификация угрозы (например, Malware, Ransomware, Spyware, Adware).
• danger_level (string/integer) — Критичность угрозы по шкале рисков (например, Low, Medium, High,).


Сервер запустится по адресу: https://localhost:6060

🗺 Эндпоинты API (Кратко)


Метод	Путь	Описание
GET	/info/{md5}	Получить расширенную информацию об угрозе по её MD5-хешу.
POST	/info	Создать новую уникальную запись об угрозе (md5, name, danger_type, danger_level).
PUT	/info/{md5}	Заменить или обновить данные существующей угрозы.


Примеры запросов

POST /info (Создание угрозы):
json
{
  "md5": "8743b52063cd84097a65d1633f5c74f5",
  "name": "WannaCry.Ransomware",
  "danger_type": "Ransomware",
  "danger_level": "Critical"
}
Используйте код с осторожностью.
GET /info/8743b52063cd84097a65d1633f5c74f5 (Успешный ответ 200):
json
{
  "md5": "8743b52063cd84097a65d1633f5c74f5",
  "name": "WannaCry.Ransomware",
  "danger_type": "Ransomware",
  "danger_level": "Critical"
}