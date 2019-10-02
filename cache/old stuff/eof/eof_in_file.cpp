#include <stdio.h>
#include <stdlib.h>


int main()
{
  char ch;
  ch = getchar();
  printf("size %lu\n", sizeof(ch));
  printf("Input Char Is :%d\n",ch & 0xff);
  printf("Input HEX Is :%x\n",ch & 0xff);
  FILE *fp;
  int c;
  fp = fopen("empty.txt","r");
  c = getc(fp) ;
  while (c!= EOF)
  {
      putchar(c);
      c = getc(fp);
      printf("%c", c);
  }
  printf("EOF got %x\n", c & 0xff);
  fclose(fp);
  getchar();
}