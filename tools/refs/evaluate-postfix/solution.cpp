#include <string>
#include <vector>
using namespace std;

int evaluate(vector<string>& tokens) {
    vector<int> stack;
    for (const string& tok : tokens) {
        if (tok.size() == 1 && string("+-*/").find(tok[0]) != string::npos) {
            int b = stack.back();
            stack.pop_back();
            int a = stack.back();
            stack.pop_back();
            switch (tok[0]) {
                case '+': stack.push_back(a + b); break;
                case '-': stack.push_back(a - b); break;
                case '*': stack.push_back(a * b); break;
                default: stack.push_back(a / b);
            }
        } else {
            stack.push_back(stoi(tok));
        }
    }
    return stack.back();
}
