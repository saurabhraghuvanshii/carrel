#include <string>
#include <vector>
using namespace std;

static void walk(int n, int open, int close, string& cur, vector<string>& out) {
    if ((int)cur.size() == 2 * n) {
        out.push_back(cur);
        return;
    }
    if (open < n) {
        cur.push_back('(');
        walk(n, open + 1, close, cur, out);
        cur.pop_back();
    }
    if (close < open) {
        cur.push_back(')');
        walk(n, open, close + 1, cur, out);
        cur.pop_back();
    }
}

vector<string> balancedStrings(int n) {
    vector<string> out;
    string cur;
    walk(n, 0, 0, cur, out);
    return out;
}
