#include <string>
#include <vector>
using namespace std;

bool isBalanced(string& text) {
    vector<char> open;
    for (char ch : text) {
        if (ch == '(' || ch == '[' || ch == '{') {
            open.push_back(ch);
            continue;
        }
        char want = ch == ')' ? '(' : ch == ']' ? '[' : '{';
        if (open.empty() || open.back() != want) return false;
        open.pop_back();
    }
    return open.empty();
}
