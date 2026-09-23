INSERT INTO teams (name, interests, skills, technologies) VALUES
    ('Data Sparks', ARRAY['analytics', 'education'], ARRAY['data analysis', 'visualization'], ARRAY['Python', 'PostgreSQL']),
    ('Green Coders', ARRAY['ecology', 'smart city'], ARRAY['backend', 'IoT'], ARRAY['Go', 'MQTT']),
    ('UX Crew', ARRAY['retail', 'accessibility'], ARRAY['product design', 'research'], ARRAY['Figma', 'Next.js']),
    ('FinTech Lab', ARRAY['finance', 'automation'], ARRAY['backend', 'security'], ARRAY['Go', 'PostgreSQL']),
    ('AI Rookies', ARRAY['AI', 'support'], ARRAY['NLP', 'frontend'], ARRAY['Python', 'TypeScript']);

INSERT INTO tasks (
    initial_description, title, topic, context, need, users, data, constraints,
    expected_result, success_criteria, contact, interaction_format, status,
    readiness_score, readiness_level, confirmed_at, published_at
) VALUES
    (
        'Нужно понять причины оттока клиентов.', 'Аналитика оттока клиентов', 'analytics',
        'Компания видит рост числа ушедших клиентов.', 'Найти основные факторы оттока.',
        'Аналитики и менеджеры продукта.', 'Обезличенная история активности за год.',
        'Нельзя передавать персональные данные.', 'Дашборд с факторами риска.',
        'Выделены минимум три проверяемых фактора.', 'product@example.test', 'Еженедельный созвон.',
        'published', 100, 'priority', NOW(), NOW()
    ),
    (
        'Нужен сервис для заявок на ремонт.', 'Учёт заявок на ремонт', 'automation',
        'Заявки сейчас принимаются в мессенджерах.', 'Собрать заявки в одном интерфейсе.',
        'Сотрудники офиса.', NULL, 'MVP за четыре недели.', 'Рабочий веб-прототип.',
        'Пользователь может создать и отследить заявку.', 'ops@example.test', 'Два созвона в неделю.',
        'published', 80, 'ready', NOW(), NOW()
    ),
    (
        'Хотим уменьшить расход электричества.', 'Мониторинг энергопотребления', 'ecology',
        'Расходы на электричество растут.', 'Показать зоны избыточного потребления.',
        'Инженеры эксплуатации.', 'Показания счётчиков по часам.', NULL,
        'Отчёт и прототип мониторинга.', NULL, 'energy@example.test', 'Асинхронно в почте.',
        'published', 75, 'ready', NOW(), NOW()
    ),
    (
        'Нужно улучшить навигацию.', 'Навигация для посетителей', 'accessibility',
        'Посетители часто не находят нужные кабинеты.', NULL, 'Посетители здания.',
        NULL, NULL, NULL, NULL, 'admin@example.test', NULL,
        'confirmed', 25, 'draft', NOW(), NULL
    ),
    (
        'Хотим чат-бота для поддержки.', NULL, 'support', NULL, NULL, NULL,
        NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 0, 'draft', NULL, NULL
    ),
    (
        'Нужно упростить адаптацию новых сотрудников.', 'Онбординг сотрудников', 'education',
        NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 0, 'draft', NULL, NULL
    ),
    (
        'Хотим понять, почему пользователи бросают оформление заказа.', 'Исследование корзины', 'retail',
        NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 0, 'draft', NULL, NULL
    ),
    (
        'Нужен удобный отчёт по обращениям клиентов.', 'Отчёт по поддержке', 'support',
        NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 0, 'draft', NULL, NULL
    ),
    (
        'Хотим автоматизировать сверку ежемесячных платежей.', 'Сверка платежей', 'finance',
        NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 0, 'draft', NULL, NULL
    );

UPDATE tasks
SET readiness_breakdown = jsonb_build_array(
    jsonb_build_object('name', 'context_and_need', 'points', CASE WHEN id = 4 THEN 10 ELSE CASE WHEN id = 5 THEN 0 ELSE 20 END END, 'max_points', 20),
    jsonb_build_object('name', 'data', 'points', CASE WHEN id IN (1, 3) THEN 20 ELSE 0 END, 'max_points', 20),
    jsonb_build_object('name', 'expected_result', 'points', CASE WHEN id IN (1, 2, 3) THEN 15 ELSE 0 END, 'max_points', 15),
    jsonb_build_object('name', 'success_criteria', 'points', CASE WHEN id IN (1, 2) THEN 15 ELSE 0 END, 'max_points', 15),
    jsonb_build_object('name', 'constraints', 'points', CASE WHEN id IN (1, 2) THEN 10 ELSE 0 END, 'max_points', 10),
    jsonb_build_object('name', 'users', 'points', CASE WHEN id IN (1, 2, 3, 4) THEN 10 ELSE 0 END, 'max_points', 10),
    jsonb_build_object('name', 'business_communication', 'points', CASE WHEN id IN (1, 2, 3) THEN 10 WHEN id = 4 THEN 5 ELSE 0 END, 'max_points', 10)
)
WHERE id <= 5;

UPDATE tasks SET
    missing_information = ARRAY['data'],
    suggestions = ARRAY['Перечислите доступные данные, примеры или источники и условия доступа.']
WHERE id = 2;

UPDATE tasks SET
    missing_information = ARRAY['constraints', 'success_criteria'],
    suggestions = ARRAY['Укажите сроки, допустимые технологии или ограничения доступа.', 'Опишите измеримые условия приёмки результата.']
WHERE id = 3;

UPDATE tasks SET
    missing_information = ARRAY['need', 'data', 'constraints', 'expected_result', 'success_criteria', 'interaction_format'],
    suggestions = ARRAY['Укажите, что бизнесу необходимо изменить в текущем процессе.', 'Перечислите доступные данные, примеры или источники и условия доступа.', 'Укажите сроки, допустимые технологии или ограничения доступа.', 'Укажите конкретный результат: например, прототип, отчёт или инструмент.', 'Опишите измеримые условия приёмки результата.', 'Опишите формат консультаций и порядок обратной связи.']
WHERE id = 4;

UPDATE tasks SET
    missing_information = ARRAY['title', 'context', 'need', 'users', 'data', 'constraints', 'expected_result', 'success_criteria', 'contact', 'interaction_format'],
    suggestions = ARRAY['Дайте задаче короткое название (не влияет на балл).', 'Опишите текущий процесс и проблему: что происходит сейчас.', 'Укажите, что бизнесу необходимо изменить в текущем процессе.', 'Назовите пользователей решения и их рабочие задачи.', 'Перечислите доступные данные, примеры или источники и условия доступа.', 'Укажите сроки, допустимые технологии или ограничения доступа.', 'Укажите конкретный результат: например, прототип, отчёт или инструмент.', 'Опишите измеримые условия приёмки результата.', 'Укажите контакт ответственного представителя бизнеса.', 'Опишите формат консультаций и порядок обратной связи.']
WHERE id = 5;

INSERT INTO proposals (task_id, team_id, solution_idea, plan, timeline, prototype_url, status) VALUES
    (1, 1, 'Модель сегментации факторов оттока.', 'Аудит данных, анализ, дашборд.', '4 недели', NULL, 'pending'),
    (1, 4, 'Сервис раннего предупреждения.', 'Подготовка признаков, правила риска, API.', '5 недель', NULL, 'accepted'),
    (2, 3, 'Простой мобильный интерфейс заявок.', 'Исследование, прототип, тестирование.', '3 недели', 'https://example.test/prototype', 'pending'),
    (2, 2, 'Веб-сервис с очередью исполнителей.', 'API, интерфейс, демонстрация.', '4 недели', NULL, 'rejected'),
    (3, 2, 'Карта потребления по зонам.', 'Импорт данных, агрегация, визуализация.', '4 недели', NULL, 'pending');
