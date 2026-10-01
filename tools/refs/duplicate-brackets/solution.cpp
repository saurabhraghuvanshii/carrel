#include <string>
using namespace std;

// For each pair, remember whether anything sat directly inside it.
bool hasPointless(string& expr) {
    string stack;
    for (char ch : expr) {
        if (ch == '(') {
            stack.push_back('(');
        } else if (ch == ')') {
            if (stack.back() == '(') return true;
            stack.pop_back();  // the '*' mark: something was inside
            stack.pop_back();  // the '('
        } else if (!stack.empty() && stack.back() == '(') {
            stack.push_back('*');
        }
    }
    return false;
}
