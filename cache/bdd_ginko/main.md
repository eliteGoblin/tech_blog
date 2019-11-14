


beforeeach: 

*  其存在context/describe的所有it都会执行
*  试做所有it都一样，不管处于nest第几位置，都会运行beforeeach的statement
*  多级beforeeach是链式关系, parent beforeeach先执行 
    ```
    beforeeach
      ...
        beforeeach
    ```
*  Describe blocks to describe the individual behaviors of your code
*  Context blocks to exercise those behaviors under different circumstances
*  Describe包含多个Context
*  BeforeEach and AfterEach blocks run for each It block
*  JustBeforeEach blocks are guaranteed to be run after all the BeforeEach blocks have run and just before the It block has run