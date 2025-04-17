ACÁ LA IDEA ES PLANTEAR UN PSEUDOCÓDIGO PARA INTERPRETAR LAS INSTRUCCIONES.

solicitarPIDyPCKernel(){
    pid, err := recibirPid(pid)

    if err != nil {
        bla bla
    }

    pc, err := recibirPC(pc)

    if err != nil {
        bla bla
    }

}

1ro. Conectarme al Kernel.

2do. Recibir PID y PC del Kernel.

3ro. Solicitar a memoria la instrucción para iniciar la ejecución, donde se ejecuta el ciclo de instrucción.

    *3.0 Traducir direcciones lógicas (del proceso) a direcciones físicas (de la memoria) [Para ésto se necesita una MMU]*

    3.1 Etapa de Fetch: Una vez recibido el PID y el PC, hay que pedirla a la memoria indicándole ese PID y ese PC.
    3.2 Etapa de Decode: Interpretar que instrucción es y si requiere una traducción.
    3.3 Etapa de Execute:
        3.3.1 NOOP: No operation, es decir, la instrucción solo va a consumir el tiempo del ciclo de instrucción.
        3.3.2 WRITE (dirección, datos): Escribe los datos de "datos" en la dirección física a partir de la dirección lógica de "dirección" (datos siempre va a ser un "string" sin espacios).
        3.3.3 READ (dirección, tamaño): Lee el valor de memoria correspondiente a la dirección física obtenida a partir de la dirección lógica de "dirección", de un tamaño determinado por el parámetro "tamaño" y lo imprime por pantalla y en el "log".
        3.3.4 GOTO (valor): Actualiza el PC del proceso al valor pasado por parámetro.

4to. Cada instrucción son llamadas Syscalls, y van a tener un nombre distinto cada una.
    4.1 IO (tiempo)
    4.2 INIT_PROC (archivo de instrucciones, tamaño)
    4.3 DUMP_MEMORY
    4.4 EXIT
    
5to. 