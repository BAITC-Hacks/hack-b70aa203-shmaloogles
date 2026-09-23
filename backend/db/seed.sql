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
        'published', 94, 'priority', NOW(), NOW()
    ),
    (
        'Нужен сервис для заявок на ремонт.', 'Учёт заявок на ремонт', 'automation',
        'Заявки сейчас принимаются в мессенджерах.', 'Собрать заявки в одном интерфейсе.',
        'Сотрудники офиса.', NULL, 'MVP за четыре недели.', 'Рабочий веб-прототип.',
        'Пользователь может создать и отследить заявку.', 'ops@example.test', 'Два созвона в неделю.',
        'published', 76, 'ready', NOW(), NOW()
    ),
    (
        'Хотим уменьшить расход электричества.', 'Мониторинг энергопотребления', 'ecology',
        'Расходы на электричество растут.', 'Показать зоны избыточного потребления.',
        'Инженеры эксплуатации.', 'Показания счётчиков по часам.', NULL,
        'Отчёт и прототип мониторинга.', NULL, 'energy@example.test', 'Асинхронно в почте.',
        'published', 58, 'workable', NOW(), NOW()
    ),
    (
        'Нужно улучшить навигацию.', 'Навигация для посетителей', 'accessibility',
        'Посетители часто не находят нужные кабинеты.', NULL, 'Посетители здания.',
        NULL, NULL, NULL, NULL, 'admin@example.test', NULL,
        'confirmed', 35, 'draft', NOW(), NULL
    ),
    (
        'Хотим чат-бота для поддержки.', NULL, 'support', NULL, NULL, NULL,
        NULL, NULL, NULL, NULL, NULL, NULL,
        'draft', 5, 'draft', NULL, NULL
    );

INSERT INTO proposals (task_id, team_id, solution_idea, plan, timeline, prototype_url, status) VALUES
    (1, 1, 'Модель сегментации факторов оттока.', 'Аудит данных, анализ, дашборд.', '4 недели', NULL, 'pending'),
    (1, 4, 'Сервис раннего предупреждения.', 'Подготовка признаков, правила риска, API.', '5 недель', NULL, 'accepted'),
    (2, 3, 'Простой мобильный интерфейс заявок.', 'Исследование, прототип, тестирование.', '3 недели', 'https://example.test/prototype', 'pending'),
    (2, 2, 'Веб-сервис с очередью исполнителей.', 'API, интерфейс, демонстрация.', '4 недели', NULL, 'rejected'),
    (3, 2, 'Карта потребления по зонам.', 'Импорт данных, агрегация, визуализация.', '4 недели', NULL, 'pending');
