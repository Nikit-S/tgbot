CREATE TABLE public.users (
        u_name varchar PRIMARY KEY,
        user_t varchar,
		l0_link varchar,
		l0_start_date date,
		l0_end_date date,
		l0_status varchar,
		
		l1_link varchar,
		l1_start_date date,
		l1_end_date date,
		l1_status varchar,
		
		l2_link varchar,
		l2_start_date date,
		l2_end_date date,
		l2_status varchar
        
);

-- Auto-generated SQL script #202204012212
INSERT INTO public.users ("u_name",user_t ,l0_link ,l0_start_date ,l0_end_date ,l0_status)
        VALUES ('Шевцов Никита Борисович', 'root', 'some link', '2021-11-26T06:22:19Z', '2021-12-26T06:22:19Z', 'сделано');
