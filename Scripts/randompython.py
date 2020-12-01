semilla =21
k = int(input("Digite numero"))
cmult=3 + 8*k
modulo = 16
generar= int(input("Cuantos quiere generar"))
generar2 = generar*4
for i in range(generar*2):
    num = (semilla*cmult)%generar2 
    print(num)
    semilla = num
