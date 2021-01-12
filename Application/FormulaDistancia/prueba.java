/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package co.com.pichincha.service.cdt.flow.bean;

import java.text.DecimalFormat;

/**
 *
 * @author julgue221
 */
public class prueba {

    public static void main(String[] args) {
        prueba prueba = new prueba();
        
        System.out.println("******Ingresando Metodo Calcular Valores*********"); 
        double s = 406.26;
        double t = -500;
        double w = -200;
        double k = 205.06;
        double h = 100;
        double p = -100;
        
        System.out.println("s = " + s); 
        System.out.println("t = " + t); 
        System.out.println("w = " + w); 
        System.out.println("k = " + k); 
        System.out.println("h = " + h); 
        System.out.println("p = " + p); 
        
        DecimalFormat df = new DecimalFormat("#"); 
        df.setMaximumFractionDigits(2);
        
        double a = prueba.calcularValorA(s, t, w, k, h, p);
        
        System.out.println("Retorno calculo a = " + df.format(a)); 
        
        double b = prueba.calcularValorB(s, t, w, k, h, p);
        
        System.out.println("Retorno calculo b = " + df.format(b)); 
        
        double c = prueba.calcularValorC(s, t, w, k, h, p);
        
        System.out.println("Retorno calculo c = " + df.format(c)); 
        
        double xp = prueba.calcularValorXPositivo(a, b, c);
        
        //double xp = -105;
        
        System.out.println("Retorno calculo x positiva = " + df.format(xp)); 
        
        double xn = prueba.calcularValorXNegativa(a, b, c);
        
        System.out.println("Retorno calculo x negativa = " + df.format(xn)); 
        
        double yp = prueba.calcularValorY(k,h,xp,p);
        
        System.out.println("Retorno calculo y positivo = " + df.format(yp)); 
        
        double yn = prueba.calcularValorY(k,h,xn,p);
        
        System.out.println("Retorno calculo y negativo = " + df.format(yn)); 
    }

    private double calcularValorA(double s, double t, double w, double k, double h, double p)
    {
        System.out.println("Empieza calculo a"); 
        
        double a = (8 * Math.pow(s, 2)) - (4 * Math.pow(t, 2)) + (8 * Math.pow(k, 2)) - (4 * Math.pow(h, 2)) - (4 * Math.pow(w, 2)) + (8 * w * p) - (4 * Math.pow(p, 2)) - (8 * t * h);
        
        System.out.println("Termina calculo a"); 
        
        return a;
    }
    
    private double calcularValorB(double s, double t, double w, double k, double h, double p)
    {
        System.out.println("Empieza calculo b"); 
        
        double b = (2 * t * Math.pow(s, 2)) - (2 * Math.pow(t, 3)) - (6 * t * Math.pow(k, 2)) - (10 * t * Math.pow(h, 2)) - (2 * t * Math.pow(w, 2)) + (4 * t * w * p) - (2 * t * Math.pow(p, 2)) - (6 * h * Math.pow(s, 2)) + (6 * h * Math.pow(t, 2)) + (2 * h * Math.pow(k, 2)) - (2 * Math.pow(h, 3)) - (2 * h * Math.pow(w, 2)) + (4 * h * w * p) - (2 * h * Math.pow(p, 2));
        
        System.out.println("Termina calculo b"); 
        
        return b;
    }
    
    private double calcularValorC(double s, double t, double w, double k, double h, double p)
    {
        System.out.println("Empieza calculo c"); 
        
        double c = (Math.pow(s, 4)) - (2 * Math.pow(s, 2) * Math.pow(t, 2)) - (2 * Math.pow(s, 2) * Math.pow(k, 2)) + (2 * Math.pow(s, 2) * Math.pow(h, 2)) - (2 * Math.pow(s, 2) * Math.pow(w, 2)) + (4 * Math.pow(s, 2) * w * p) - (2 * Math.pow(s, 2) * Math.pow(p, 2)) + (Math.pow(t, 4)) + (2 * Math.pow(t, 2) * Math.pow(k, 2)) - (2 * Math.pow(t, 2) * Math.pow(h, 2)) + (2 * Math.pow(t, 2) * Math.pow(w, 2)) - (4 * Math.pow(t, 2) * w * p) + (2 * Math.pow(t, 2) * Math.pow(p, 2)) + (Math.pow(k, 4)) - (2 * Math.pow(k, 2) * Math.pow(h, 2)) - (2 * Math.pow(k, 2) * Math.pow(w, 2)) + (4 * Math.pow(k, 2) * w * p) - (2 * Math.pow(k, 2) * Math.pow(p, 2)) + (Math.pow(h, 4)) + (2 * Math.pow(h, 2) * Math.pow(w, 2)) - (4 * Math.pow(h, 2) * w * p) + (2 * Math.pow(h, 2) * Math.pow(p, 2)) + (Math.pow(w, 4)) - (4 * Math.pow(w, 3) * p) + (4 * Math.pow(w, 2) * Math.pow(p, 2)) - (4 * w * Math.pow(p, 3)) + (Math.pow(p, 4));
        
        System.out.println("Termina calculo c"); 
         
        return c;
    }
    
    private double calcularValorXPositivo(double a, double b, double c)
    {
        System.out.println("Empieza calculo x positiva"); 
        
        double raiz = Math.pow(b, 2) - (4 * a *c);
        
        double x = ((-b) + (Math.sqrt(raiz)))/(2*a);
        
        System.out.println("Termina calculo x positiva"); 
        
        return x;
    }
    
    private double calcularValorXNegativa(double a, double b, double c)
    {
        System.out.println("Empieza calculo x negativa"); 
        
        double raiz = Math.pow(b, 2) - (4 * a * c);
        
        double x = ((-b) - Math.sqrt(raiz))/(2*a);
        
        System.out.println("Termina calculo x negativa"); 
        
        return x;
    }
    
    private double calcularValorY(double k, double h, double x, double p)
    {
        System.out.println("Empieza calculo y"); 
        
        double hx = h-x;
        
        double raiz = Math.pow(k, 2) - Math.pow((hx), 2);
        
        double y = (-Math.sqrt(raiz) + p);
        
        System.out.println("Termina calculo y"); 
        
        return y;
    }
}
