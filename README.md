# HighloadEmail #

## Статистика из открытых источников ##

Исходя из данных [wikipedia](https://ru.wikipedia.org/wiki/Население_России) в России живет порядка 146 млн человек. [yandex.radar](https://radar.yandex.ru/yandex?month=2020-08) же говорит, что число уникальных пользователей Яндекс почты составляет около 17млн человек в месяц , в процентом соотношении это 11% населения страны. Сылаясь на [similarweb](https://www.similarweb.com/website/mail.yandex.ru/) у почты Яндекса общее количество посещений в месяц примерно 400млн(У почты Mail.ru - более 540млн), со средним временем нахождения пользователя на сайте - 8 минут. Общее количество активных аккаунтов почты mail.ru составляет [100млн](https://corp.mail.ru/ru/company/portal/) 
Лимит на объем одного письма без вложения сделаем равным 10MB. Есть возможность прикреплять файлы(вложения),максимальный объем вложений в сумме - 20MB

## Запросы ##

   -  Регистрация
   -  Авторизация
   -  Постраничное получение списка последних писем
   -  Отправка и пересылка писем
   -  Удаление писем
   -  Чтение писем
   -  Скачивание вложений из писем
   
## Средний пользователь ##

Предположим, что средний пользователь отправляет и получает по 5 писем в день(по статистике яндекса эти цифры в 2 раза меньше), а список писем получает 10 раз в день. Ссылаясь на [источник_1](https://habr.com/ru/company/yandex/blog/199432/) и [источник2](https://habr.com/ru/post/321756/) около 10% писем содержат вложения. При этом средний размер письма без вложения 4KB, а с ним 400КВ. Средний размер почтового ящика -  6000(высчитано эмпирически) писем, из которых 2000 были получены за последние 2 года(в терминологии S3 хранилищ будем считать, что это те данные, которые относятся к HotBox). Суточная аудитория Яндекс почты порядка [6.2 млн человек](https://radar.yandex.ru/yandex?month=2020-08), при этом каждый, в среднем остается на сайте около 8 минут(для простоты вычислений возьмем 5).

## Эксперимент ##

Вспоминая нашего усредненного пользователя, проведем эксперимент. Откроем в браузере вкладку network и смоделируем поведения нашего пользователя. получаем 500 запросов, сгенерированных пользователем.

## Нагрузка ##

Нагрузка общая, исходя из результатов эксперимента
       
    RPS = (6.2*10^6 * 500)/(24*60*60) ~ 36000
    
### Целевые цифры проекта ###

5 млн пользователей в день

80 млн активных аккаунтов (54% населения РФ, 80% [Mail.ru](https://corp.mail.ru/ru/))

     RPS = 30000(5000 - бэкенд, 20000 фронтенд/веб-сервер, 5000 - вспомогательные службы) запросов в секунду
     
Примем во внимание, что пиковое время для почты это около 10 утра по мск, в это время сервисом пользуется более 60% от основной активной аудитории, то есть 

      5000000 * 0.6= 3000000 пользователей

1) Получения списка последних писем (50 штук).Данные будут браться из отдельной таблицы коротких сообщений. Согласно исследованию, данный запрос будет весить примерно 23KB данных,без учета иконок почт отправителей( 1 и 2).Учитывая цифры для среднего пользователя из пункта 2, получаем:

       Постраничное получение
       
       RPS: 5 000 000 * 10 / 24 * 60 * 60 ~ 600
       
       Передаваемый трафик: 5 000 000 * 10 * 23KB / 24 * 60 * 60 ~ 13800KB ~ 14MB
       
       Пиковый 3000000 * 23КБ = 526ГБит/с
       
Получается, что для обработки получения списка писем хватит мощности одного сервера.
2) Чтение одного письма - согласно исследованию, средний вес одного письма(в конечной для клиента форме)

12KB без учета прикрепляемого файла.(3 и 4 ) Примем, что пользователь читает хотя бы 5 писем в день. Тогда получим:

       Чтение письма 
       
       RPS: 5 000 000 * 5 / 24 * 60 * 60 ~ 290
       
       Передаваемый трафик: 5 000 000 * 5 * 12KB / 24 * 60 * 60 ~ 3480KB ~ 3,3MB
       
       Пиковый 3000000 * 12КБ = 137ГБит/с
       
3) Отправка одного письма - согласно исследованию, средний вес одного письма без вложения 4кб

Примем, что пользователь отправляет хотя бы 5 писем в день. Тогда получим:

       Отправка письма 
       
       RPS: 5 000 000 * 5 / 24 * 60 * 60 ~ 290
       
       Передаваемый трафик: 5 000 000 * 5 * 4KB / 24 * 60 * 60 ~ 1157KB ~ 1,1MB
       
       Пиковый 3000000 * 4КБ * 5 = 457ГБит/с
       
Для данной задачи мощности все того же одного сервера с LevelDB вполне достаточно.
4) Написание письма - согласно пункту 2 мы считаем, что средний вес письма без вложения- 4KB, а с вложением - 400KB.

       5 000 000 пользователей * 5 писем в день = 25 000 000 писем в день всего
       25 000 000 * 0,9 = 22 500 000 писем без вложений
       
       22 500 000 * 4Кб = 90 000 000Кб весят все письма за день =~ 85Гб
       
       RPS писем без вложения: 5 000 000 * 5 / 24 * 60 * 60 = 289 
       
       Пиковый 3000000 * 5 = 114ГБит/с


       25 000 000 * 0,1 = 2 500 000 писем с вложениями
       
       2 500 000 * 400Кб =~ 953Гб весят письма с вложениями
       
       Вес вложений составляет 99% веса письма, тогда примем , что мы получаем 940Гб вложений в день.
       
       RPS писем c вложениями: 5 000 000 * 1 / 24 * 60 * 60 = 58 
       
       Пиковый 3000000 * 1 = 22ГБит/с

## Логическа схема базы данных ##

![image](https://user-images.githubusercontent.com/47527934/115372266-620b2100-a1d3-11eb-8054-ef3b80d34aeb.png)

## Физическая схема базы данных ##

### User ###

id     | name          | surname      | birthday        | gender   | accountname | password    | phone
------ | ------------- | ------------ | --------------- | -------- | ----------- | ----------- | ------
bigint | varchar(20)   | varchar(20)  | timestamp:date  | boolean  | varchar(30) | varchar(32) | bigint

Итог: 126 байт на одного пользователя

### Letter ###

id     | sender  | receiver        | theme        | text              | date           | isRead  | folderId        | answerOn  | fileId          |
------ | ------- | --------------- | ------------ | ----------------- | -------------- | ------- | --------------- | --------- | --------------- |
bigint | bigint  | bigint array[5] | varchar(100) | varchar(11000000) | timestamp:date | boolean | bigint array[5] | bigint    | bigint array[5] |
                        

Итог: 11000252 байт на одного пользователя

### File ###

id     | name          | extension  | size | path
------ | ------------- | ---------- | ---- | -----------
bigint | varchar(20)   | varchar(8) | int  | varchar(20)

Итог: 60 байт на одного пользователя

### Folder ###

id     | name        | letterCount | userId
------ | ----------- | ----------- | -------
bigint | varchar(20) | int         | bigint

Итог: 40 байт на одного пользователя

В качестве СУБД будет использоваться LevelDB. Для каждого пользователя на сервере с СУБД будет создана папка, в которой будут храниться все связанные с ним сущности(Иконки,письма,вложения,данные о папках).Это позволит получать сообщения пользователя, конкретно из его стораджа, а не из какой-то общей таблицы,в которой будет миллионы сообщений, что очень замедлит поиск.

Таблица со списком пользователей будет храниться в KeyValue хранилище, в нашем случае это будет redis. Это позволит быстро находить получателей письма и искать людей, с которыми пользователь переписывался.

Согласно пунктам 7 и 9 из benchmark При размере данных в 100B и 128MB кеша:

Random Reads - 137,231 ops/sec
Random Writes - 176,929 ops/sec
Batch Random Writes - 229,095 entries/sec

Данные загружаются на сервер со скоростью записи жесткого диска, наша база записывает в себя по сути только снипетт на сохраняемый файл. 


Прикрепленные файлы будем хранить в хранилище S3, все файлы будут находиться в Cold боксах, так как согласно эмперическим иследованиям пользователи не так часто смотрят вложения писем, соответственно "холодное" хранение куда выгоднее "горячего".


      Итоговый размер базы: 80_000_000(аккаунтов) * 6000(размер ящика) * 126(байт на пользователя) * 11000252(байт письмо) = 6.6529524096*10^20 байт
      
       

## Технологический стек ## 

Бэкенд: Главным языком для серверной разработки будет выбран Golang, который позволяет быстро и эффективно писать код.В его основе лежит теоретическая модель CSP.

Scalling by adding more of the same

Golang UK Conference 2017 by Arne Claus
Она обеспечит хорошую масштабируемость Бэкенда, а также позволит избежать блокировок при конкурентном доступе. В основном код будет написан с использованием обширной стандартной библиотеки, роутером будет выбран fasthttp, который показывает наилучшие показатели производительности. Безусловно , нужно учитывать трудности со сборщиком мусора, который будет снижать производительность раз за какое-то △t времени.Поэтому необходимо будет учитывать не только среднюю производительность , но и ее нижнюю границу.

Фронтенд: Выбор CSS,HTML,TypeScript/JavaScript , как технологий для фронтенд составляющей проекта будет обусловлен их популярностью и изученностью годами.Для JavaScript существует уже огромное количество библиотек ,которые как и упрощают верстку, так и способны влиять на производительность кода.TypeScript способен и ускорить разработку из-за отлавливания ошибок при компиляции, и обезопасить наш будущий код от непредсказуемого поведения.

Протоколы: В данном проекте подразумевается использование https и http2 для связи клиента с бэкендом.Также будет реализована связь с помощью websocket протокола, для отображения новых пришедших пользователю сообщений(как пример) в режиме реального времени. Связь микросервисов на бэкенде будет осуществляться классическим для Golang'а способом - через gRPC+Protobuf, что является одним из самых эффективных способов общения как с точки зрения скорости, так и с точки зрения размера передаваемого трафика.

Еще раз распишем путь нашего запроса по шагам.Для примера возьмем написания письма с вложением. Пускай пользователь уже зарегистрирован.Рассматривать будем только позитивные сценарии.
1) Пользователь запрашивает наш сайт. Веб-сервер(о нем речь пойдет ниже) отдает ему необходимую статику,а потом перенаправляет запросы с заголовками на бэкенд.
2) Бэкенд просит Redis отдать ему id пользователя и сервер с его данными по cookie.Далее бэкенд сервер осуществляет соединение с БД по полученному адресу и запрашивает оттуда базовую информацию о пользователе.После получения отдает все обратно клиенту.
3) У клиента отобразилась стартовая страница со всеми данными.Он переходит в раздел написания письма.Создает его,прикрепляет файл и совершает отправку.
4) Отправка с вложениями по REST API методу отличается от отправки без вложения.Веб-сервер
отправляет запрос на бэкенд.
5) Бэкенд сообщений узнает у Redis'a в какой сервер БД направляется это сообщение и отправляет по этому адресу. 7) Сервер БД получает письма, сохраняет их в папку пользователя, записывает необходимую информацию об этом в levelDB.

Физическая схема БД

![physical](https://user-images.githubusercontent.com/47527934/118355947-890f0580-b57b-11eb-9ab9-68e52f3fd20a.png)




## Расчет нагрузки и оборудования ## 


БД:
В 4 пункте были рассмотрены многие цифры, касательно системы писем и пользователей. Предположим, что у одного юзера в среднем 6000(2000 "горячих") писем в ящике всего(отправленные и входящие во всех папках) и все папки. Получаем

Общий объем для хранения писем без вложений = 80млн * 5600 * 4KB ~ 1668 TB
Общий объем для хранения писем с вложениями= 80млн * 600 * 400KB ~ 17881 TB

Общий "горячий" объем для хранения писем без вложений= 80млн * 1800 * 4KB ~ 537 TB
Общий "холодный" объем для хранения писем без вложений= 1668TB - 537TB ~  1131 TB

Общий "горячий" объем для хранения писем с вложениями= 80млн * 200 * 400KB ~ 5960 TB
Общий "холодный" объем для хранения писем с вложениями= 17881TB - 5960TB ~  11921 TB

Средний вес одного "горячего" пользователя со всеми данными  = 
0,2KB(данные о пользователе) + (1800 * 4KB)(письма без вложений) + (200 * 400KB)(с вложениями) 
+ (6 * 0,3KB)(папки) = 88000 KB ~ 79 MB


Средний вес одного "холодного" пользователя без посторонних данных  = 
(3600 * 4KB)(письма без вложений) + (400 * 400KB)(с вложениями) = 174 400 KB ~ 170MB
Для хранения данных об одном "горячем" пользователе нужно 79MB.

Выходит, что если мы закупим сервера вида:

CPU(cores) |	RAM(GB) |	SSD(GB)/2-unit(10)
--------|-----------|------
8 |	32 |	4096 x 10


то на один из них с учетом запаса под рост данных в 50% на нем уместится 20480GB/79MB ~ 265 000 пользователей. Всего таких серверов понадобится 80 000 000 / 265 000 = 300 штук Надежность этих данных требуется высокая, поэтому нам придется делать по 2 реплики к каждому. Итого 600 реплик, 900 серверов всего.

Что касается архивных данных , то тут мы воспользуемся конфигурацией(Опробуем новинки 2020 года в лице 18 - терабайтных жестких дисков)

CPU(cores) |	RAM(GB)	| HDD(TB)/2-unit(10)
--------|-----------|------
8	| 32 |	18 * 10


На такой сервер с учетом запаса данных в 66% на нем уместится 90TB/170MB ~ 555 000 пользователей. Всего таких серверов понадобится 80 000 000 / 555 000 = 145 штук На холодные сервера не сильная нагрузка и не такие высокие требования к надежности, поэтому нам хватит по 1 реплики к каждому. Итого еще 145 реплик , 290 серверов всего.

Сервер Redis:
Из 30000 запросов примерно в половине требуется информация о пользователе. Примем RPS = 15 000. Согласно redis_bench мы видим, что 8 ядер нам более чем достаточно.Возьмем

CPU(cores) |	RAM(GB) |	SSD(GB)
--------|-----------|------
8 |	32 |	512 


Однако надежность этого сервиса крайне важна, поэтому сделаем ему целых 3 реплики.Итого 4 сервера. У Redis'а есть удобная встроенная репликация, что является плюсом.

Фронтенд:
Основная нагрузка на фронтенд сервера - парсинг HTML файлов.Для этого воспользуемся библиотекой на языке Go - exp/html.Учитвая bench, можно посчитать , что на 32-ядерном сервере будет 3000 распаршенных документов в секунду.Если предположить, что все наши 15000 запросов на фронтенд будут требовать парсинг HTML,то можем получить итоговое количество серверов.

Много физической памяти фронтенд серверу не особо надо,она будет потрачена в основном на кеширование запросов и хранения документов, а вот RAM ему понадобиться побольше.Возьмем:

CPU(cores) |	RAM(GB) |	SSD(GB)
--------|-----------|------
32|	64|	512


В количестве **5 штук **.Так же добавим каждому по 2 запасных.Итого 15 серверов

Балансировщик:

В пункте 2 мы выбрали общую нагрузку в 30000RPS. Будем использовать nginx.Согласно nginx_test при среднем объеме требуемого пользователю контента в 10MB , мы видим , что можем обойтись 16 ядрами и 32GB RAM. Объем документной составляющей примерно 5MB на человека.Тогда подсчитаем, что 5 * 5 000 000(количество пользователей в день) = 25 000 000 MB/день ~ 300MB/s Но учтем , что это при распределении на 24 часа.На деле значение может быть выше.

CPU(cores) |	RAM(GB) |	SSD(GB)
--------|-----------|------
16 |	64	| 512


Возьмем 10 серверов.Но оставлять их без поддержки - опасно.Возьмем еще по 2(Итого 30) запасных дополнительно,которые будут подменять первый при падении.Между ними настроим CARP. При падении главного веб-сервера его IP подхватит другой и полностью заменит его.

Бэкенд:

Бэкенд должен соответствовать производительности БД, но благодаря микросервисной архитектуре мы грамотно разбиваем нагрузку. Согласно 5 и 6, мы видим что большая часть запросов - получение картинок(50/50 статика - пользовательские) и вспомогательные запросы для клиентской стороны.Из суммарных 30000RPS на бэкенд приходится лишь примерно 5000, которые и распределяются по микросервисам. Если мы возьмем 5 микросервисов, 2 для авторизации и 3 для работы с сообщениями,то каждому из них будет достаточно следующих параметров:

CPU(cores) |	RAM(GB) |	SSD(GB)
--------|-----------|------
8	| 16	| 512


Каждому из 5 микросервисов нужен по 1 запасному.Итого 10 серверов. Вспомним про перенос данных из "горячих" серверов в "холодные".Тогда будет 6 серверов для работы с сообщениями, к каждому по запасному.Итого 8 серверов основных. 16 серверов всего.

### Итог ###

Назначение     | количество        | Оборудование 
-------------- | ----------------- | ----------- 
Микросервисы         | 10       | CPU(cores) 8	RAM(GB) 16	SSD(GB) 512
Работа с сообщениями | 6       | CPU(cores) 8	RAM(GB) 16	SSD(GB) 512
Балансировщик        | 10 * 3      | CPU(cores) 8	RAM(GB) 64	SSD(GB) 512
Фронтенд         | (1 + 2) * 5      | CPU(cores) 32	RAM(GB) 16	SSD(GB) 512
Redis         | 1 + 3      | CPU(cores) 8	RAM(GB) 32	SSD(GB) 512
Холодные данные         | 145 * 2      | CPU(cores) 8	RAM(GB) 32	HDD(TB)/2-unit(10) 18 * 10
Горячие данные         | 300 * 3      | CPU(cores) 8	RAM(GB) 32	SSD(GB)/2-unit(10) 4096 x 10

		     

## Хостинг и расположение серверов ##


Для наших объемов данных экономически не выгодно использовать облачные сервисы, поэтому мы реализовали все хранение сами.Наш сервис ориентирован на Россию,но все же не стоит забыть и о СНГ.Эффективнее всего расположить основные сервера в западной части России, а расположение архивных можно и в средней/восточной ее части. Для уменьшения задержек можно добавить еще один фронтенд сервер в средней/восточной части России с его двумя запасами.Тогда количество фронтенд-серверов примем равным 30.

## Балансировка нагрузки ##

![Balancer](https://user-images.githubusercontent.com/47527934/116789470-a772fc80-aab7-11eb-9067-48dd3146aab7.png)

![CARP](https://user-images.githubusercontent.com/47527934/116789474-b0fc6480-aab7-11eb-8727-1ec9297f4eab.png)

Выше уже говорилось, что будем использовать nginx для балансировки нагрузки на 7 уровне модели OSI(L7). А так как их у нас больше 1 экземпляра, то еще будем использовать DNS балансировку.Алгоритм балансировки DNS и некоторые его тонкости представлены на картинке выше.По умолчанию клиент будет пытаться отправить запрос на первый из IP адресов в списке.DNS сервер будет перемешивать(сдвиг на один) массив IP адресов.На всякий случай укажем небольшой TTL(1 минуту),чтоб у нас все равно оставалось больше гибкости при работе с IP адресами наших веб-серверов. Алгоритм балансировки на уровне L7 аналогичен(Клиент = веб-сервер, в нем IP адреса бэкендов, выбирается один и в него запрос идет).


## Отказоустойчивость ##

Для обеспечения устойчивости нашим многочисленным серверам БД необходимо поддерживать у них реплики.Сделать это только встроенным средствами LevelDB не выйдет. Будем использовать подход логической репликации. Все изменения в базе данных происходят в результате вызовов её API – например, в результате выполнения SQL-запросов. Очень заманчивой кажется идея выполнять одну и ту же последовательность запросов на двух разных базах. Для репликации необходимо придерживаться двух правил:
 
1. Нельзя начинать транзакцию, пока не завершены все транзакции, которые должны закончиться раньше. Так на рисунке ниже нельзя запускать транзакцию D, пока не завершены транзакции A и B.
2. Нельзя завершать транзакцию, пока не начаты все транзакции, которые должны закончиться до завершения текущей транзакции. Так на рисунке ниже даже если транзакция B выполнилась мгновенно, завершить её можно только после того, как начнётся транзакция C.

Обычно для логической репликации используют детерминированные запросы. Детерминированность запроса обеспечивается двумя свойствами:

1. запрос обновляет (или вставляет, или удаляет) единственную запись, идентифицируя её по первичному (или уникальному) ключу;
2. все параметры запроса явно заданы в самом запросе.

База-реплика открыта и доступна не только на чтение, но и на запись. Это позволяет использовать реплику для выполнения части запросов, в том числе для построения отчётов, требующих создания дополнительных таблиц или индексов.

Логическая репликация предоставляет ряд возможностей, отсутствующих в других видах репликации:

1. настройка набора реплицируемых данных на уровне таблиц (при физической репликации – на уровне файлов и табличных пространств, при блочной репликации – на уровне томов);
2. построение сложных топологий репликации – например, консолидация нескольких баз в одной или двунаправленная репликация;
3. уменьшение объёма передаваемых данных;
4. репликация между разными версиями СУБД или даже между СУБД разных производителей;
5. обработка данных при репликации, в том числе изменение структуры, обогащение, сохранение истории.
