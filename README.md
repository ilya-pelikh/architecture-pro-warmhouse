### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи (IoT) могут создавать, читать, редактировать, удалять информацию о сенсорах
- Система поддерживает возможность создавать, читать, редактировать, удалять информацию о сенсорах

**Мониторинг температуры:**

- Пользователи могут запросить температуру в конкретной локации.
- Система поддерживает возможность запросить температуру в конкретной локации.

### 2. Анализ архитектуры монолитного приложения

- Язык програмирования: GO
- База данных PostgreSQl
- В качестве маршрутизатора приложения используется GIN
- Приложение использует слоистую архитектуру:
	- handlers
	- models
	- services
	- bd

### 3. Определение доменов и границы контекстов

Домены:
- Управление сенсорами
- Температура

### **4. Проблемы монолитного решения**

- Сложность раздельного развития доменов
- Повышение связанности модулей, в следствие чего повышение нагрузки при поддержке и развитии

### 5. Визуализация контекста системы — диаграмма С4

[![C1. Диаграмма контекстов](https://www.plantuml.com/plantuml/png/JKz1IiKm4Dtd54Ft5Ro0VD2DQ_K08N-GGjf2aWhkh6vTkXTl42p5MDkUON8ZPwA216OcxxrvZpdue5qOMgj24op2-Ua9q0ibYZJb1wuhlmYqq4vRVgIPbZnoZqff4te75RqozPMVwOFxEqKHoRy2xU76erCEJT3ThKKMlr4g-_EFxIgsb1ZOEItnF36SH-OthWoAQR6wklS1MsLiFnoXkkolw_maLaNoAPEa5-a5tZK6lL8I3tLzzhVyI_o73jWT-Vnl)](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/schemas/current/c1.puml)


# Задание 2. Проектирование микросервисной архитектуры

**Диаграмма контейнеров (Containers)**

[![C2. Диаграмма контейнеров](https://www.plantuml.com/plantuml/png/XLH1QzHG4Bw_Np7alXuyUf2gLIg8h6R5Sve-PA5D97cAHOHkfNfG4A5dY_yWhbtCRfFqBtpl7pap8HZJq5BocfdvvfjlvarsAtTHvjuSx1itT_kM7jo-HJX5YGA__H6V5VyLX0qQrFuOAwnHOyMX1ajZejYb6GKuO_F2VyKVSQCGeDW1HJczo6rcJ6PZ5oMlk4I0doZzGZ3AxJ_YKaijSS6kceFVn5bnmI-e2ETesE0deFU6A4uFNgU9JHm5Zx8qSLf4tpZOB_C3l_oHY56YDMcP1Xk_G3dXGxfOogaCB-tkD4lV7VOhg5Q5YM9HifBsiEPiRbJpjktRUMn3cyc_wvBRdXjsbFup1nROyDk9i1QeHQC_KSQTUSvJFRoPaWa_wSGA8HLjcX_6DecoUVgFO2zoqMwKor2-u_IwFhRlZy1xob7s9OwXRn35AH-CaPsWoUHkUszqXtPcisber59_TniSUlumdD6zwMFJ-GGm67i_r8aRZGrdCx2lKhC0dBNhhfFVe22KX2MlrALiFMhfLZfjlv9AR_FS_CSoeak9bN8b6y0hKL3xtV1JSU2udVS8KNTLJF4q3PyzCGFcNh4UQApLpgSm0MtPP4z-3Vu2)](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/schemas/next/c2.puml)

**Диаграмма компонентов (Components)**

[![C3. Диаграмма компонентов](https://www.plantuml.com/plantuml/png/VLJDRXCn4BxxAKQzXmit3gYGzC2jIa2Sjnkhh9JkLXiNKI6aRFZX00A18d51LBp14bgnJKdo2kDNu4duF8s95naKoHdRttpVUEPBTiScASt9qKYPp8mqTM_U89vnmpL_jnEiF_6fYHUyGh36kkGBNjefeZdzcio8l_85a7DkR4Msig_SFTSuQX3DbIfrQfBy4EGx_8M_ELQaB-qRlYHUOxlYXQsuwM0H9hVkK3AcX_fWlooUdAgqE3ekSfqRK8HoVEuaKEKbSb6GPFNquopL0rswhvRy_qj5zbD556WDKfqdPLPqYzoKnMYaIn7ORkFmZRHWDt4Sk9SVPlarxZhy1r_wI0zFlFgVNHpSryV5E3D5UIhexLOzKsgOcEGW6RjMVE1Pj0N8yLIzkqhv4L2dSpjuEDsY_jxWFjtPlwU2Ku3yYu7ezX_qY1ko4_lSleRCoaty1zUyna752-Mno0igHIecjR8lnO3YDrTTf5omjmYT1qBIjXSSxvduJIm6k3Uso9otCcKWal9P469wuAxiM-_TralyLFvA_hJ5aFlvaDY_GWWFYcnZtod1nzfm1kqhwLwy2VQzDFW5hXQE3tSuSw0fDbCHhZrzj24VyPQrKmaQtMfAK108MkROrVdG_NNy1W00)](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/schemas/next/c3.puml)

**Диаграмма кода (Code)**

[![C4. Диаграмма кода](https://www.plantuml.com/plantuml/png/VLDHZjCm4FtFAKRzwWAfXCHNXAg5xKyG0eco1x3YsMB9iOFjHaI0j7i7Bi0DgeML2a9m1UT6Z3VPjL8NLAcUcVSoyzx4ERME6RUjPUQ4amXdJnRWl_Yzt_bl_INz_wBV3lnF_zj_fs3k_r1vJyMR_h80_pMS_jhl-oiWmWs5E_-ZloBwzKcMLP9P2ojiHOLdMZcZfKG37pE0Xr9ODo1sVbB0IqgV3ldgxUjbO6nm9Cm59yOz63236te0pKWm-7wBrXsoy0DMMqUZwvefJcqUGMIsMl3iyrJa6aquRXNc1XjjXTEckslxvhxoVmV_bdz4Exyz9eoMIS5P53vGIf9V3PKx-Qd7b1gT2aLJBrY5YPtXbLm8nPz9cPePeETToG-zMDQc3k9I3yQKTICZ5ZSjU07dv-Kgn8hLMC3Q6Q4s8NTTCyajEaUXBU25rUgLvXYgeraxAX-tiKBIBPpGQXb59LZGaCvjcA6AGsFNrJkiMPmNMYLtEi-NqpkKJRA0fvVes7vWiVYKKeztYsDWzdW63_8SPkIjcwKcirDKF7peVm40)](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/schemas/next/c4.puml)

# Задание 3. Разработка ER-диаграммы
[![ER. Сервис сценариев](https://www.plantuml.com/plantuml/png/fLDHQkf05FtFAWO_XRW04L7YF0ZlgG_rFVedOt9hfoGpOUQcQBNGHNGXGd_wqKMaE-fDj2In0McBOCAvvoRdd9mPM1G6qpXoK684ZCywsLFsaXsovxT7-d-bty5nGD5uoreZ_j_t-Bp5X6KXtCa0BDivZBMP35cNfIbTUhs_l_jz0bstc2WH0yqi6gdM1OBRv0onW4Ztbj1RgzMg00CTnqA59VXG2fbxV3AS-TEZ51k04aRgxuYXpLK4-S9AwmY48jG6Dn2cX2u5qW1b39HJdE1zauQbDnsFXvDHNSlo9F-YfdAZIlbrfBHFCTTci44mIaJCm3KOK04Gjv9OtM0Y39ccZNgwYROIfLO4IeMm1jFeZLznxz_2dqxEtC46WZJV_wa_2XPJMvTSbB8fVPPRJ2B02vzcmV_mCL_C_5DN4IA8WLvmkQzV2jn0e4tOh9IU951KorlMTBzttVtkSrVEA7hN5BZpSV8wdPgnRYczGcaeP9Ku01Ncn_yT)](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/schemas/next/er.puml)


# Задание 4. Создание и документирование API

### 1. Тип API

В качестве типа API для первичной реализации я бы выбрал REST API, из-за его простоты и легкости интеграции, поддержки.
Как вариант для будущего апдейта, можно рассматривать AsyncAPI, интегрировав RabbitMQ/Kafka и переведя систему в Event driven design.

### 2. Документация API
[OpenAPI-спецификация](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/open-api.yml)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/devices` | Возвращает список зарегистрированных девайсов, их настройки и поддерживаемые команды. |
| `POST` | `/devices` | Регистрирует новый девайс вместе с его настройками и доступными командами. |
| `GET` | `/devices/{deviceId}` | Возвращает конкретный девайс по UUID. |
| `GET` | `/scenarios` | Возвращает список созданных сценариев с командами и расписанием. |
| `POST` | `/scenarios` | Создаёт сценарий из упорядоченного набора команд для девайсов. |
| `GET` | `/scenarios/{scenarioId}` | Возвращает конкретный сценарий по UUID. |
| `PUT` | `/scenarios/{scenarioId}` | Полностью заменяет настройки, расписание и команды сценария. |
| `DELETE` | `/scenarios/{scenarioId}` | Удаляет сценарий. |
| `POST` | `/scenarios/{scenarioId}/executions` | Запускает сценарий немедленно и возвращает информацию о запуске. |
| `GET` | `/scenario-executions/{executionId}` | Возвращает состояние запуска: `pending`, `running`, `succeeded` или `failed`. |

# Задание 5. Работа с docker и docker-compose

[smart_home/Dockerfile](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/apps/smart_home/Dockerfile)  
[temperature/Dockerfile](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/apps/temperature/Dockerfile)  
[apps/docker-compose](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/apps/docker-compose.yml)  

```
cd apps
docker compose up --build -d
```

# **Задание 6. Разработка MVP**

[apps/device_service](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/apps/device_service)  

[apps/history_service](https://github.com/ilya-pelikh/architecture-pro-warmhouse/blob/warmhouse/apps/history_service)  

NOTE: не стал разрабатывать сервис сценариев, тк на мой взгляд он не является необходимым для MVP
Логику взаимодействия новой архитектуры можно увидеть в С3 диаграмме