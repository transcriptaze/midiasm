qwerty(X) :- print(X).

uiop :-
	open('qwerty.txt', read, S),
    read_file(S),
    close(S).

read_file(S) :-
    get_char(S,X),
    read_file(S,X).

read_file(_,end_of_file) :-
	write('---\n').

read_file(S,C) :-
	write(C),
    get_char(S,X),
    read_file(S,X).
