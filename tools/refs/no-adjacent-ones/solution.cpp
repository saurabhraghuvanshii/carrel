#include <string>
#include <vector>
using namespace std;

static void build(int n, string& cur, vector<string>& out) {
    if ((int)cur.size() == n) {
        out.push_back(cur);
        return;
    }
    cur.push_back('0');
    build(n, cur, out);
    cur.pop_back();
    if (cur.empty() || cur.back() != '1') {
        cur.push_back('1');
        build(n, cur, out);
        cur.pop_back();
    }
}

vector<string> noAdjacentOnes(int n) {
    vector<string> out;
    string cur;
    build(n, cur, out);
    return out;
}
